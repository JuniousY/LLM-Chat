package llm

type Chunk struct {
	ID   int64     `milvus:"name:id"`
	V    []float32 `milvus:"name:vector"`
	Text string    `milvus:"name:text"`
}
