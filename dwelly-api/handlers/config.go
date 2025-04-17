package handlers

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/koredeycode/dwelly/internal/database"
	"github.com/redis/go-redis/v9"
)

type APIConfig struct {
	DB         *database.Queries
	Redis      *redis.Client
	Cloudinary *cloudinary.Cloudinary
	Validate   *validator.Validate
}
