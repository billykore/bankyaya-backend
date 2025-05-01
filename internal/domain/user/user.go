// Package user provides domain logic related to user operations.
//
// It includes business rules for user authentication, registration, profile management,
// and other identity-related behaviors. This package defines core entities and interfaces
// for handling user data and workflows, ensuring separation of concerns from infrastructure
// (e.g., database, HTTP) and application layers.
package user

import "time"

// User represents a user in the system.
// It contains personal and account-related information, including identifiers,
// contact details, and device information.
type User struct {
	ID            int
	CIF           string
	Password      string
	AccountNumber string
	FullName      string
	Email         string
	PhoneNumber   string
	NIK           string
	Device        *Device
}

// The Device represents a user's device information used for authentication,
// push notifications, and security checks.
type Device struct {
	FirebaseId    string
	DeviceId      string
	IsBlacklisted bool
}

// Valid checks whether the provided device credentials match the device's credentials.
func (d *Device) Valid(firebaseId string, deviceId string) bool {
	return d.ValidFirebaseId(firebaseId) && d.ValidDeviceId(deviceId)
}

// ValidFirebaseId checks whether the provided Firebase ID matches the device's Firebase ID.
func (d *Device) ValidFirebaseId(firebaseId string) bool {
	return firebaseId == d.FirebaseId
}

// ValidDeviceId checks whether the provided device ID matches the device's ID.
func (d *Device) ValidDeviceId(deviceId string) bool {
	return deviceId == d.DeviceId
}

// The Token represents an authentication token issued to a user or client.
// It contains the access token string, its type (e.g., "Bearer"), and the expiration time.
// This struct is typically used in authentication flows such as OAuth2 or JWT-based systems.
type Token struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
}
