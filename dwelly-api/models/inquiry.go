package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

type Inquiry struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ListingID uuid.UUID `json:"listing_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Status    string    `json:"status"`
}

func DatabaseInquirytoInquiry(dbInquiry database.Inquiry) Inquiry {
	return Inquiry{
		ID:        dbInquiry.ID,
		CreatedAt: dbInquiry.CreatedAt,
		UpdatedAt: dbInquiry.UpdatedAt,
		ListingID: dbInquiry.ListingID,
		SenderID:  dbInquiry.SenderID,
		Status:    dbInquiry.Status,
	}
}

func DatabaseInquiriestoInquiries(dbInquirys []database.Inquiry) []Inquiry {
	inquiries := make([]Inquiry, len(dbInquirys))
	for i, inquiry := range dbInquirys {
		inquiries[i] = DatabaseInquirytoInquiry(inquiry)
	}
	return inquiries

}
