package models

type ChatRequest struct {
	AppId          int     `json:"app_id"`
	ConversationId int     `json:"conversation_id"`
	Msg            *string `json:"msg"`
	DocumentId     *string `json:"document_id"`
}
