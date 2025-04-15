package database

import (
	"time"

	"github.com/google/uuid"
)

// For Listing
func (l Listing) GetID() uuid.UUID        { return l.ID }
func (l Listing) GetCreatedAt() time.Time { return l.CreatedAt }
func (l Listing) GetUpdatedAt() time.Time { return l.UpdatedAt }
func (l Listing) GetUserID() uuid.UUID    { return l.UserID }
func (l Listing) GetIntent() string       { return l.Intent }
func (l Listing) GetTitle() string        { return l.Title }
func (l Listing) GetDescription() string  { return l.Description }
func (l Listing) GetPrice() string        { return l.Price }
func (l Listing) GetLocation() string     { return l.Location }
func (l Listing) GetCategory() string     { return l.Category }
func (l Listing) GetStatus() string       { return l.Status }

// For GetListingByIDRow
func (l GetListingByIDRow) GetID() uuid.UUID        { return l.ID }
func (l GetListingByIDRow) GetCreatedAt() time.Time { return l.CreatedAt }
func (l GetListingByIDRow) GetUpdatedAt() time.Time { return l.UpdatedAt }
func (l GetListingByIDRow) GetUserID() uuid.UUID    { return l.UserID }
func (l GetListingByIDRow) GetIntent() string       { return l.Intent }
func (l GetListingByIDRow) GetTitle() string        { return l.Title }
func (l GetListingByIDRow) GetDescription() string  { return l.Description }
func (l GetListingByIDRow) GetPrice() string        { return l.Price }
func (l GetListingByIDRow) GetLocation() string     { return l.Location }
func (l GetListingByIDRow) GetCategory() string     { return l.Category }
func (l GetListingByIDRow) GetStatus() string       { return l.Status }

// For ListAllListingsRow

func (l ListAllListingsRow) GetID() uuid.UUID        { return l.ID }
func (l ListAllListingsRow) GetCreatedAt() time.Time { return l.CreatedAt }
func (l ListAllListingsRow) GetUpdatedAt() time.Time { return l.UpdatedAt }
func (l ListAllListingsRow) GetUserID() uuid.UUID    { return l.UserID }
func (l ListAllListingsRow) GetIntent() string       { return l.Intent }
func (l ListAllListingsRow) GetTitle() string        { return l.Title }
func (l ListAllListingsRow) GetDescription() string  { return l.Description }
func (l ListAllListingsRow) GetPrice() string        { return l.Price }
func (l ListAllListingsRow) GetLocation() string     { return l.Location }
func (l ListAllListingsRow) GetCategory() string     { return l.Category }
func (l ListAllListingsRow) GetStatus() string       { return l.Status }

// For ListUserListingsRow
func (l ListUserListingsRow) GetID() uuid.UUID        { return l.ID }
func (l ListUserListingsRow) GetCreatedAt() time.Time { return l.CreatedAt }
func (l ListUserListingsRow) GetUpdatedAt() time.Time { return l.UpdatedAt }
func (l ListUserListingsRow) GetUserID() uuid.UUID    { return l.UserID }
func (l ListUserListingsRow) GetIntent() string       { return l.Intent }
func (l ListUserListingsRow) GetTitle() string        { return l.Title }
func (l ListUserListingsRow) GetDescription() string  { return l.Description }
func (l ListUserListingsRow) GetPrice() string        { return l.Price }
func (l ListUserListingsRow) GetLocation() string     { return l.Location }
func (l ListUserListingsRow) GetCategory() string     { return l.Category }
func (l ListUserListingsRow) GetStatus() string       { return l.Status }

// For SearchListingsRow
func (l SearchListingsRow) GetID() uuid.UUID        { return l.ID }
func (l SearchListingsRow) GetCreatedAt() time.Time { return l.CreatedAt }
func (l SearchListingsRow) GetUpdatedAt() time.Time { return l.UpdatedAt }
func (l SearchListingsRow) GetUserID() uuid.UUID    { return l.UserID }
func (l SearchListingsRow) GetIntent() string       { return l.Intent }
func (l SearchListingsRow) GetTitle() string        { return l.Title }
func (l SearchListingsRow) GetDescription() string  { return l.Description }
func (l SearchListingsRow) GetPrice() string        { return l.Price }
func (l SearchListingsRow) GetLocation() string     { return l.Location }
func (l SearchListingsRow) GetCategory() string     { return l.Category }
func (l SearchListingsRow) GetStatus() string       { return l.Status }

// For Inquiry
func (i Inquiry) GetID() uuid.UUID        { return i.ID }
func (i Inquiry) GetCreatedAt() time.Time { return i.CreatedAt }
func (i Inquiry) GetUpdatedAt() time.Time { return i.UpdatedAt }
func (i Inquiry) GetListingID() uuid.UUID { return i.ListingID }
func (i Inquiry) GetSenderID() uuid.UUID  { return i.SenderID }
func (i Inquiry) GetStatus() string       { return i.Status }

// func (i Inquiry) GetMessages() []Message  { return i.Messages }

// For GetInquiryByIDWithMessagesRow
func (i GetInquiryByIDWithMessagesRow) GetID() uuid.UUID        { return i.ID }
func (i GetInquiryByIDWithMessagesRow) GetCreatedAt() time.Time { return i.CreatedAt }
func (i GetInquiryByIDWithMessagesRow) GetUpdatedAt() time.Time { return i.UpdatedAt }
func (i GetInquiryByIDWithMessagesRow) GetListingID() uuid.UUID { return i.ListingID }
func (i GetInquiryByIDWithMessagesRow) GetSenderID() uuid.UUID  { return i.SenderID }
func (i GetInquiryByIDWithMessagesRow) GetStatus() string       { return i.Status }
