package llm

type ChatCompletionChunk struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Delta        Delta       `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason interface{} `json:"finish_reason"`
}

type Delta struct {
	Content          *string `json:"content"`
	ReasoningContent *string `json:"reasoning_content"`
}

type Usage struct {
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

func (chunk ChatCompletionChunk) GetContent() string {
	if len(chunk.Choices) == 0 {
		return ""
	}
	if chunk.Choices[0].Delta.Content == nil {
		return ""
	}
	return *chunk.Choices[0].Delta.Content
}

func (chunk ChatCompletionChunk) GetReasoning() string {
	if len(chunk.Choices) == 0 {
		return ""
	}
	if chunk.Choices[0].Delta.ReasoningContent == nil {
		return ""
	}
	return *chunk.Choices[0].Delta.ReasoningContent
}
