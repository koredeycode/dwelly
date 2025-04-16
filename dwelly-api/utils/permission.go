package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

func IsListingOwner(db *database.Queries, r *http.Request, listingId uuid.UUID, userId uuid.UUID) (bool, error) {
	listing, err := db.GetListingByID(r.Context(), listingId)
	if err != nil {
		return false, err
	}

	if listing.UserID == userId {
		return true, nil
	}

	return false, nil

}

func IsInquirySender(db *database.Queries, r *http.Request, inquiryId, userId uuid.UUID) (bool, error) {
	inquiry, err := db.GetInquiryById(r.Context(), inquiryId)
	if err != nil {
		return false, err
	}

	if inquiry.SenderID == userId {
		return true, nil
	}
	return false, nil
}

func IsMessageSender(db *database.Queries, r *http.Request, messageId, userId uuid.UUID) (bool, error) {
	message, err := db.GetMessage(r.Context(), messageId)
	if err != nil {
		return false, err
	}

	if message.SenderID == userId {
		return true, nil
	}
	return false, nil
}
