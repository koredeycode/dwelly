package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GetUUIDParam(uuidStr, obj string) (uuid.UUID, string) {
	if uuidStr == "" {
		return uuid.Nil, fmt.Sprintf("missing %s ID", obj)
	}
	id, err := uuid.Parse(uuidStr)
	if err != nil {

		return uuid.Nil, fmt.Sprintf("invalid %s Id", obj)
	}
	return id, ""
}
