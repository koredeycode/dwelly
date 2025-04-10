package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/koredeycode/dwelly/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

func DatabaseUsertoUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
	}
}

func DatabaseUserstoUsers(dbUsers []database.User) []User {
	users := make([]User, len(dbUsers))
	for i, user := range dbUsers {
		users[i] = DatabaseUsertoUser(user)
	}
	return users

}
