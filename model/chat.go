package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChatMessage struct {
	gorm.Model
	ChatID   uint   `gorm:"index"`
	SenderID uint   `gorm:"index"`
	Sender   string `gorm:"size:255"`
	Content  string `gorm:"size:1000"`
}

type Chat struct {
	gorm.Model
	ChatID     uuid.UUID `gorm:"type:char(36);unique;not null"`
	SenderID   int       `gorm:"index"`
	ReceiverID int       `gorm:"index"`
}

type Message struct {
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	ChatID    string    `json:"chat_id"`
	SenderID  uint      `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}
