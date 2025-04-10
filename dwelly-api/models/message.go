package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Content   string    `json:"content"`
	InquiryID uuid.UUID `json:"inquiry_id"`
	SenderID  uuid.UUID `json:"sender_id"`
}

func DatabaseMessagetoMessage(dbMessage database.Message) Message {
	return Message{
		ID:        dbMessage.ID,
		CreatedAt: dbMessage.CreatedAt,
		UpdatedAt: dbMessage.UpdatedAt,
		Content:   dbMessage.Content,
		InquiryID: dbMessage.InquiryID,
		SenderID:  dbMessage.SenderID,
	}
}

func DatabaseMessagestoMessages(dbMessages []database.Message) []Message {
	messages := make([]Message, len(dbMessages))
	for i, message := range dbMessages {
		messages[i] = DatabaseMessagetoMessage(message)
	}
	return messages

}
