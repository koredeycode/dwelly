package models

import (
	"time"

	"github.com/google/uuid"
)

type InquiryLike interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetListingID() uuid.UUID
	GetSenderID() uuid.UUID
	GetStatus() string
}

type Inquiry struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`
	ListingID uuid.UUID `json:"listing_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Status    string    `json:"status"`

	Messages []Message `json:"messages"`
}

func DatabaseInquiryToInquiry[T InquiryLike](i T) Inquiry {
	return Inquiry{
		ID:        i.GetID(),
		CreatedAt: i.GetCreatedAt(),
		// UpdatedAt: i.GetUpdatedAt(),
		ListingID: i.GetListingID(),
		SenderID:  i.GetSenderID(),
		Status:    i.GetStatus(),
	}
}

func DatabaseInquiriesToInquiries[T InquiryLike](is []T) []Inquiry {
	inquiries := make([]Inquiry, len(is))
	for i, inquiry := range is {
		inquiries[i] = DatabaseInquiryToInquiry(inquiry)
	}
	return inquiries

}
