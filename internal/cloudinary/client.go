// internal/cloudinary/client.go
package cloudinaryutil

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewClient() (*cloudinary.Cloudinary, error) {

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Printf("Failed to initialize Cloudinary: %v", err)
		return nil, err
	}

	return cld, nil
}
