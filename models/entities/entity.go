package entities

import (
	"gorm.io/datatypes"
	"time"
)

type Conversations struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID        int       `gorm:"not null" json:"app_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Introduction string    `gorm:"type:text" json:"introduction,omitempty"`
	ModelName    string    `gorm:"type:varchar(50)" json:"model_name,omitempty"`
	Status       string    `gorm:"type:varchar(255);default:'0';not null" json:"status"`
	FromUserID   string    `gorm:"type:varchar(255);not null" json:"from_user_id"`
	IsDeleted    bool      `gorm:"type:tinyint(1);default:0;not null" json:"is_deleted"`
	CreatedAt    time.Time `gorm:"type:date;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:date;not null" json:"updated_at"`
}

type Messages struct {
	ID                      int            `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID                   int            `json:"app_id"`
	ModelProvider           string         `gorm:"type:varchar(255)" json:"model_provider,omitempty"`
	ModelID                 string         `gorm:"type:varchar(255)" json:"model_id,omitempty"`
	OverrideModelConfigs    string         `gorm:"type:text" json:"override_model_configs,omitempty"`
	ConversationID          int            `json:"conversation_id"`
	Inputs                  datatypes.JSON `gorm:"type:json;not null" json:"inputs"`
	Query                   string         `gorm:"type:text;not null" json:"query"`
	Message                 datatypes.JSON `gorm:"type:json;not null" json:"message"`
	MessageTokens           int            `gorm:"default:0;not null" json:"message_tokens"`
	MessageUnitPrice        float64        `gorm:"type:numeric(10,4);not null" json:"message_unit_price"`
	Answer                  string         `gorm:"type:text;not null" json:"answer"`
	AnswerTokens            int            `gorm:"default:0;not null" json:"answer_tokens"`
	AnswerUnitPrice         float64        `gorm:"type:numeric(10,4);not null" json:"answer_unit_price"`
	ProviderResponseLatency float64        `gorm:"default:0;not null" json:"provider_response_latency"`
	TotalPrice              float64        `gorm:"type:numeric(10,7)" json:"total_price,omitempty"`
	Currency                string         `gorm:"type:varchar(255);not null" json:"currency"`
	FromSource              string         `gorm:"type:varchar(255);not null" json:"from_source"`
	FromUserID              string         `gorm:"type:varchar(255);not null" json:"from_user_id"`
	CreatedAt               time.Time      `gorm:"type:date;not null" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"type:date;not null" json:"updated_at"`
	Status                  string         `gorm:"type:varchar(255);default:'0';not null" json:"status"`
	Error                   string         `gorm:"type:text" json:"error,omitempty"`
}
