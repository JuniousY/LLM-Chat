package service

import (
	"LLM-Chat/config"
	"LLM-Chat/models"
	"LLM-Chat/models/entities"
	"LLM-Chat/models/llm"
	"LLM-Chat/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"gorm.io/datatypes"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type llmMsg struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

const contextPrompt = "Use the following context as your learned knowledge, inside <context></context> XML tags.\n\n<context>\n{{#context#}}\n</context>\n\nWhen answer to user:\n- If you don't know, just say that you don't know.\n- If you don't know when you are not sure, ask for clarification.\nAvoid mentioning that you obtained the information from the context.\nAnd answer according to the language of the user's question.\n\n"

func ChatCompletion(c *gin.Context) {
	writer := c.Writer
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "text/event-stream")
	writeDate := func(data interface{}) {
		writer.Write([]byte("data: "))
		writer.Write(utils.Marshal(data))
		writer.Write([]byte("\n\n"))
		writer.Flush()
	}
	writerDoneCh := make(chan interface{})
	writerDoneClose := new(int32)

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 详细错误处理
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	sysPrompt := "You are a helpful assistant"
	if req.DocumentId != nil {
		// mock
		cnt, promoptContext := loadChatContext(c, req.Msg, req.DocumentId)
		if cnt > 0 {
			params := map[string]string{
				"context": promoptContext,
			}
			sysPrompt = sysPrompt + "\n\n" + utils.ReplacePlaceholders(contextPrompt, params)
		}
	}
	chatReq := llm.NewChatRequest()
	msgs := []llm.Message{
		{
			Content: sysPrompt,
			Role:    "system",
		},
	}

	// 从db读取历史记录
	var preMsgData *entities.Messages
	var messageStr string

	config.DB.Raw("SELECT * FROM messages WHERE conversation_id = ?", req.ConversationId).Scan(&preMsgData)
	if preMsgData != nil {
		log.Println(preMsgData)
		preMsgs := utils.Unmarshal[[]llmMsg]([]byte(preMsgData.Message))
		log.Println(preMsgs)
		messageStr = utils.MarshalString(preMsgs)
		log.Println("messageStr-" + messageStr)
		for _, preMsg := range preMsgs {
			msgs = append(msgs, llm.Message{
				Content: preMsg.Text,
				Role:    preMsg.Role,
			})
		}
		msgs = append(msgs, llm.Message{Content: preMsgData.Answer, Role: "assistant"})
	}

	// 追加用户输入
	msgs = append(msgs, llm.Message{
		Content: *req.Msg,
		Role:    "user",
	})

	chatReq.Messages = msgs
	conversationContext := msgs[1:]
	response, err := utils.RequestDeepSeek(chatReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var chatChunks []llm.ChatCompletionChunk
	go func() {
		for response.HasNext() {
			v := response.Read()
			chatChunks = append(chatChunks, v)
			writeDate(v)
		}
		saveMessage(req, conversationContext, chatChunks)
		if atomic.CompareAndSwapInt32(writerDoneClose, 0, 1) {
			close(writerDoneCh)
		}
	}()

	timer := time.NewTimer(time.Minute * 5)
	defer timer.Stop()

	debugChunk(chatChunks)
	select {
	case <-writer.CloseNotify():
		response.Close()
		return
	case <-writerDoneCh:
		return
	case <-timer.C:
		response.Close()
		log.Println("timeout")
		return
	}
}

func debugChunk(chatChunks []llm.ChatCompletionChunk) {
	for _, chatChunk := range chatChunks {

		// 输出结果（自动处理空字符串拼接）
		fmt.Println(chatChunk.GetContent() + " " + chatChunk.GetReasoning())
	}
}

// 简化版本
func saveMessage(chatReq models.ChatRequest, messages []llm.Message, chunks []llm.ChatCompletionChunk) {
	var answer strings.Builder
	writeReason := true
	for _, chunk := range chunks {
		if chunk.GetReasoning() != "" {
			answer.WriteString(chunk.GetReasoning())
		}
		if writeReason && chunk.GetContent() != "" {
			answer.WriteString("\n")
			writeReason = false
		}
		if chunk.GetContent() != "" {
			answer.WriteString(chunk.GetContent())
		}
	}

	msg := entities.Messages{
		AppID:                   chatReq.AppId,
		ModelProvider:           "",
		ModelID:                 "",
		OverrideModelConfigs:    "",
		ConversationID:          chatReq.ConversationId,
		Inputs:                  datatypes.JSON("{}"),
		Query:                   *chatReq.Msg,
		Message:                 utils.StructToDatatypesJSON(messages),
		MessageTokens:           0,
		MessageUnitPrice:        0,
		Answer:                  answer.String(),
		AnswerTokens:            0,
		AnswerUnitPrice:         0,
		ProviderResponseLatency: 0,
		TotalPrice:              0,
		Currency:                "CNY",
		FromSource:              "",
		FromUserID:              "",
		CreatedAt:               time.Time{},
		UpdatedAt:               time.Time{},
		Status:                  "",
		Error:                   "",
	}
	log.Println("save" + utils.MarshalString(msg))
	config.DB.Save(&msg)
}

func loadChatContext(ctx *gin.Context, msg *string, documentId *string) (int, string) {
	log.Println("loadChatContext:" + *documentId)

	cli := config.MilvusCli

	// mock
	embeddings, _ := GenEmbedding([]string{*msg})
	embedding := embeddings[0].Embedding
	const collectionName = "my_rag_collection"

	searchOpt := milvusclient.NewSearchOption(
		collectionName, // collectionName
		5,              // limit
		[]entity.Vector{entity.FloatVector(embedding)},
	).
		WithOutputFields("id", "vector", "text")
	resultSets, err := cli.Search(ctx, searchOpt)
	if err != nil {
		log.Fatal("failed to perform basic ANN search collection: ", err.Error())
	}

	var contextStr bytes.Buffer
	cnt := 0
	for _, resultSet := range resultSets {
		var data []*llm.Chunk
		err := resultSet.Fields.Unmarshal(&data)
		if err != nil {
			log.Fatal(err)
			return cnt, contextStr.String()
		}
		for _, chunk := range data {
			contextStr.WriteString(chunk.Text)
			cnt++
		}
	}
	if cnt > 0 {
		return cnt, contextStr.String()
	} else {
		return cnt, ""
	}
}
