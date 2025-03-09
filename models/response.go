package models

type ChatChunk struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskID         string `json:"task_id"`
	Answer         string `json:"answer"`
}
