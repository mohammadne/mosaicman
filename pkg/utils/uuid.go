package utils

import (
	"strings"

	uuid "github.com/google/uuid"
)

func NewUUID(key string) string {
	sha1 := uuid.NewSHA1(uuid.New(), []byte(key))
	return strings.Replace(sha1.String(), "-", "", -1)
}
