package repository

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
)

// UserRepository defines methods for managing user persistence.
type UserRepository interface {
	// GetUserDataByPhoneNumber retrieves a user from the database by their phone number.
	// Requires a context and a string phone number as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error)

	// GetDeviceById retrieves a device from the database by its unique ID.
	// Requires a context and a strings device ID as input parameters.
	// Returns a Device object and an error if retrieval fails.
	GetDeviceById(ctx context.Context, deviceId string) (*entity.Device, error)
}
