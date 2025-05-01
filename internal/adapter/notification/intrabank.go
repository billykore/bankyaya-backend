package notification

import (
	"context"

	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"go.bankyaya.org/app/backend/internal/pkg/notification/firebase"
)

type IntrabankNotification struct {
	firebase *firebase.Client
}

func NewIntrabankNotification(firebaseClient *firebase.Client) *IntrabankNotification {
	return &IntrabankNotification{
		firebase: firebaseClient,
	}
}

func (n *IntrabankNotification) Notify(ctx context.Context, notification *intrabank.Notification) error {
	err := n.firebase.Send(ctx, &firebase.Message{
		FirebaseId: notification.FirebaseId,
		Title:      notification.Subject,
		Body:       notification.String(),
	})
	if err != nil {
		return err
	}
	return nil
}
