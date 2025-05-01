package codes

// Code represents an integer type typically used for error or status codes.
type Code int

const (
	// OK is success code.
	_ Code = iota // OK

	// BadRequest represents a code indicating a bad request error.
	BadRequest

	// Unauthenticated represents a code indicating an unauthenticated error.
	Unauthenticated

	// Forbidden represents a code indicating a forbidden error.
	Forbidden

	// NotFound represents a code indicating a resource could not be found.
	NotFound

	// Conflict represents a code indicating a conflict error.
	Conflict

	// Internal represents a code indicating an internal server error.
	Internal
)
