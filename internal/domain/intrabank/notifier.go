package intrabank

import "context"

// Notifier sends intrabank transfer notifications to users.
type Notifier interface {
	// Notify sends a transfer notification to the specified user.
	Notify(ctx context.Context, notification *Notification) error
}
