package test

import (
	"LLM-Chat/service"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file ", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGenEmbedding(t *testing.T) {
	input := []string{"花椰菜又称菜花、花菜，是一种常见的蔬菜。"}
	embeddings, e := service.GenEmbedding(input)
	if e != nil {
		t.Fatal(e)
	}

	res := embeddings[0].Embedding
	embeddingDim := len(res)
	fmt.Println(embeddingDim)

	fmt.Println(res)
}
