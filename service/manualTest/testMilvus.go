package main

import (
	"LLM-Chat/service"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const collectionName = "my_rag_collection"

type Record struct {
	A    int64     `milvus:"name:id"`
	V    []float32 `milvus:"name:vector"`
	Text string    `milvus:"name:text"`
}

// 参考 https://milvus.io/docs/zh/build-rag-with-milvus.md
// doubao dimension 2560
func ManualTestInsert() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli := initCli(ctx)
	if cli == nil {
		return
	}

	textLines := getTestTextLines()

	has, err := cli.HasCollection(ctx, milvusclient.NewHasCollectionOption(collectionName))
	if err != nil {
		fmt.Println(err)
		return
	}
	if has {
		cli.DropCollection(ctx, milvusclient.NewDropCollectionOption(collectionName))
	}

	dim := 2560
	createOpt := milvusclient.SimpleCreateCollectionOptions(collectionName, int64(dim))
	createOpt.WithMetricType(entity.IP)
	createOpt.WithConsistencyLevel(entity.ClStrong)
	cli.CreateCollection(ctx, createOpt)

	embeddings, _ := service.GenEmbedding(textLines)

	ids := make([]int64, 0, len(textLines))
	vectors := make([][]float32, 0, len(textLines))
	for i := range textLines {
		ids = append(ids, int64(i))
		vectors = append(vectors, embeddings[i].Embedding)
	}
	insertOpt := milvusclient.NewColumnBasedInsertOption(collectionName).
		//WithInt64Column("id", ids).
		WithFloatVectorColumn("vector", dim, vectors).
		WithVarcharColumn("text", textLines)

	insertResult, err := cli.Insert(ctx, insertOpt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert result:", insertResult)
}

func initCli(ctx context.Context) *milvusclient.Client {
	milvusAddr := os.Getenv("MILVUS_HOST")
	cli, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: milvusAddr,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cli
}

func getTestTextLines() []string {
	// 存储分割后的文本行
	var textLines []string
	// 定义要遍历的目录和文件模式
	pattern := "resources/milvus_docs_en/faq/*.md"
	// 遍历匹配的文件
	err := filepath.WalkDir("resources/milvus_docs_en/faq/", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 检查文件是否匹配模式
		match, err := filepath.Match(pattern, path)
		if err != nil {
			return err
		}
		if match && !d.IsDir() {
			// 读取文件内容
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			// 将文件内容按 # 分割
			lines := strings.Split(string(content), "# ")
			// 将分割后的行添加到结果列表
			textLines = append(textLines, lines...)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	}
	for _, line := range textLines {
		fmt.Println(line)
	}
	return textLines
}

func ManualTestQuery1() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli := initCli(ctx)
	if cli == nil {
		return
	}

	collection, err := cli.DescribeCollection(ctx, milvusclient.NewDescribeCollectionOption(collectionName))
	if err != nil {
		// handle error
	}
	fmt.Println(collection)

	rs, err := cli.Get(ctx, milvusclient.NewQueryOption(collectionName).
		WithIDs(column.NewColumnInt64("id", []int64{0, 1, 2, 3})))
	if err != nil {
		// handle error
	}
	fmt.Println(rs.GetColumn("id"))
}

func ManualTestQuery() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli := initCli(ctx)
	if cli == nil {
		return
	}

	question := "How is data stored in milvus?"
	embeddings, _ := service.GenEmbedding([]string{question})
	embedding := embeddings[0].Embedding

	searchOpt := milvusclient.NewSearchOption(
		collectionName, // collectionName
		3,              // limit
		[]entity.Vector{entity.FloatVector(embedding)},
	).
		WithOutputFields("id", "vector", "text")
	resultSets, err := cli.Search(ctx, searchOpt)
	if err != nil {
		log.Fatal("failed to perform basic ANN search collection: ", err.Error())
	}
	for _, resultSet := range resultSets {
		log.Println("IDs: ", resultSet.IDs)
		log.Println("Fields: ", resultSet.Fields)
		var data []*Record
		err := resultSet.Fields.Unmarshal(&data)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, record := range data {
			fmt.Println("String array:", record.Text)
		}

		log.Println("Scores: ", resultSet.Scores)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file ", err)
		os.Exit(1)
	}

	//ManualTestInsert()
	ManualTestQuery()
}
