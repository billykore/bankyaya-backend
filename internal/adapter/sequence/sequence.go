package sequence

import "github.com/google/uuid"

type UUID struct{}

func New() *UUID {
	return &UUID{}
}

func (u *UUID) Generate() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
