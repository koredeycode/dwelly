package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

type Listing struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ListingID   uuid.UUID `json:"user_id"`
	Intent      string    `json:"intent"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`

	Images []string `json:"images"`
}

func DatabaseListingtoListing(dbListing database.Listing) Listing {

	return Listing{
		ID:          dbListing.ID,
		CreatedAt:   dbListing.CreatedAt,
		UpdatedAt:   dbListing.UpdatedAt,
		Intent:      dbListing.Intent,
		Title:       dbListing.Title,
		Description: dbListing.Description,
		Price:       dbListing.Price,
		Location:    dbListing.Location,
		Category:    dbListing.Category,
		Status:      dbListing.Status,

		Images: []string{}, // Assuming images are not included in the dbListing struct
	}
}

func DatabaseListingstoListings(dbListings []database.Listing) []Listing {
	listings := make([]Listing, len(dbListings))
	for i, listing := range dbListings {
		listings[i] = DatabaseListingtoListing(listing)
	}
	return listings

}
