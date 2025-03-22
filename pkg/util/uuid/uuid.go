package uuid

import "github.com/google/uuid"

// New generates a new UUID (Version 7) and returns its string representation or an error.
func New() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
