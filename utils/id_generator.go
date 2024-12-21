package utils

import "github.com/google/uuid"

// I wanted to use uuid to generate unique IDs
func GenerateID() string {
	return uuid.New().String()
}
