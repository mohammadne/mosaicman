package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
