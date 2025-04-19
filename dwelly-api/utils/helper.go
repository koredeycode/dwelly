package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GetUUIDParam(uuidStr, obj string) (uuid.UUID, string) {
	if uuidStr == "" {
		return uuid.Nil, fmt.Sprintf("Missing %s ID", obj)
	}
	id, err := uuid.Parse(uuidStr)
	if err != nil {

		return uuid.Nil, fmt.Sprintf("Invalid %s Id", obj)
	}
	return id, ""
}

// ExtractValidationErrors converts validator.ValidationErrors into a slice of error strings
func ExtractValidationErrors(err error) []string {
	var ve validator.ValidationErrors
	var errorMessages []string

	if errors.As(err, &ve) {
		for _, e := range ve {
			errorMessages = append(errorMessages,
				fmt.Sprintf("Field '%s' failed on the '%s' rule", e.Field(), e.Tag()),
			)
		}
	} else {
		// If not a ValidationErrors type, include the general error
		errorMessages = append(errorMessages, err.Error())
	}

	return errorMessages
}
