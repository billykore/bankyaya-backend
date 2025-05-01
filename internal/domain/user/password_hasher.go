package user

// PasswordHasher defines an interface for hashing and verifying passwords.
type PasswordHasher interface {
	// Hash generates a hashed representation of the given password.
	// password: The plain-text password to hash.
	// Returns the hashed password string and an error if hashing fails.
	Hash(password string) (string, error)

	// Compare checks whether the provided plain-text password matches the given hashed password.
	// password: The plain-text password input. hashed: The stored hashed password.
	// Returns true if the password matches, otherwise false.
	Compare(password, hashed string) bool
}
