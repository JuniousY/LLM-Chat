package service

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GenEmbedding(input []string) ([]openai.Embedding, error) {
	config := openai.DefaultConfig(os.Getenv("ARK_API_KEY"))
	config.BaseURL = os.Getenv("ARK_BASE_URL")
	client := openai.NewClientWithConfig(config)
	model := os.Getenv("ARK_MODEL_EMBEDDING")

	//fmt.Println("----- embeddings request -----")
	req := openai.EmbeddingRequestStrings{
		Input:          input,
		Model:          openai.EmbeddingModel(model),
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
	}

	resp, err := client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		fmt.Printf("embeddings error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}
