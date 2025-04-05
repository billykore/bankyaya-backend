package sequence

import "github.com/google/uuid"

type Sequence struct{}

func New() *Sequence {
	return &Sequence{}
}

func (seq *Sequence) Generate() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
