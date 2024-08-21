package models

import (
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GetCurrentFormatedTime() string {
	return time.Now().Format(time.RFC3339)
}
