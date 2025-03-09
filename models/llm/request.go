package llm

// Message 结构体
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// ResponseFormat 结构体
type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type Tool struct {
	// todo
}

// 请求结构体
type ChatRequest struct {
	Messages         []Message      `json:"messages"`
	Model            string         `json:"model"` // deepseek-chat deepseek-reasoner
	FrequencyPenalty float64        `json:"frequency_penalty"`
	MaxTokens        int            `json:"max_tokens"`
	PresencePenalty  float64        `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	Stop             []string       `json:"stop"`
	Stream           bool           `json:"stream"`
	StreamOptions    *StreamOptions `json:"stream_options"`
	Temperature      float64        `json:"temperature"`
	TopP             float64        `json:"top_p"`
	Tools            *Tool          `json:"tools"`
	ToolChoice       string         `json:"tool_choice"`
	Logprobs         bool           `json:"logprobs"`
	TopLogprobs      *int           `json:"top_logprobs"`
}

func NewChatRequest() ChatRequest {
	return ChatRequest{
		Model:            "deepseek-reasoner",
		FrequencyPenalty: 0.0,
		MaxTokens:        2048,
		PresencePenalty:  0.0,
		ResponseFormat:   ResponseFormat{Type: "text"},
		Stop:             nil,
		Stream:           true,
		StreamOptions:    &StreamOptions{IncludeUsage: true},
		Temperature:      1.0,
		TopP:             1.0,
		Tools:            nil,
		ToolChoice:       "none",
		Logprobs:         false,
		TopLogprobs:      nil,
	}
}
