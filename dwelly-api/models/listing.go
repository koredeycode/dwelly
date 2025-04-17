package models

import (
	"time"

	"github.com/google/uuid"
)

type ListingLike interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetUserID() uuid.UUID
	GetIntent() string
	GetTitle() string
	GetDescription() string
	GetPrice() string
	GetLocation() string
	GetCategory() string
	GetStatus() string
}

// implement the rest...

type Listing struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
	Intent      string `json:"intent"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Location    string `json:"location"`
	Category    string `json:"category"`
	Status      string `json:"status"`

	Images []string `json:"images"`
}

func DatabaseListingToListing[T ListingLike](l T) Listing {

	return Listing{
		ID:        l.GetID(),
		UserID:    l.GetUserID(),
		CreatedAt: l.GetCreatedAt(),
		// UpdatedAt:   l.GetUpdatedAt(),
		Intent:      l.GetIntent(),
		Title:       l.GetTitle(),
		Description: l.GetDescription(),
		Price:       l.GetPrice(),
		Location:    l.GetLocation(),
		Category:    l.GetCategory(),
		Status:      l.GetStatus(),

		Images: []string{}, // Assuming images are not included in the dbListing struct
	}
}

func DatabaseListingsToListings[T ListingLike](ls []T) []Listing {
	listings := make([]Listing, len(ls))
	for i, listing := range ls {
		listings[i] = DatabaseListingToListing(listing)
	}
	return listings

}
