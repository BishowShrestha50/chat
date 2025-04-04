package service

import (
	"chat/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IChatService interface {
	CreateChat(senderID int, receiverID int, chatID uuid.UUID) error
	SaveMessage(message *model.ChatMessage) error
	GetChatBySenderAndReceiverID(user1ID int, user2ID int) (*model.Chat, error)
	GetChatHistoryBetweenUsers(user1ID, user2ID int) (*[]model.ChatMessage, error)
}
type ChatService struct {
	DB *gorm.DB
}

func NewChatService(db *gorm.DB) IChatService {
	return &ChatService{DB: db}
}
func (c *ChatService) GetChatBySenderAndReceiverID(user1ID int, user2ID int) (*model.Chat, error) {
	var chat *model.Chat
	err := c.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *ChatService) SaveMessage(message *model.ChatMessage) error {
	return c.DB.Create(message).Error
}

func (c *ChatService) CreateChat(senderID int, receiverID int, chatID uuid.UUID) error {
	chat := &model.Chat{
		ChatID:     chatID,
		SenderID:   senderID,
		ReceiverID: receiverID,
	}
	return c.DB.Create(chat).Error
}

func (c *ChatService) GetChatHistoryBetweenUsers(user1ID, user2ID int) (*[]model.ChatMessage, error) {
	var messages *[]model.ChatMessage
	err := c.DB.Table("chat_messages").
		Select("chat_messages.*").
		Joins("inner join chats on chats.id = chat_messages.chat_id").
		Where("(chats.sender_id = ? AND chats.receiver_id = ?) OR (chats.sender_id = ? AND chats.receiver_id = ?)", user1ID, user2ID, user2ID, user1ID).
		Order("chat_messages.created_at asc").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
