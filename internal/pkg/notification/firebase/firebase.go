package firebase

import (
	"context"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
)

const timeout = 10 * time.Second

type Message struct {
	FirebaseID string
	Title      string
	Body       string
}

type Client struct {
	fcmClient *messaging.Client
}

func New() *Client {
	return newClient()
}

func newClient() *Client {
	log := logger.New()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	app, err := firebase.NewApp(ctx, &firebase.Config{})
	if err != nil {
		log.Errorf("failed to initialize firebase app: %v", err)
		return nil
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		log.Errorf("failed to initialize firebase messaging client: %v", err)
		return nil
	}

	return &Client{
		fcmClient: fcmClient,
	}
}

func (c *Client) Send(ctx context.Context, message *Message) error {
	_, err := c.fcmClient.Send(ctx, &messaging.Message{
		Token: message.FirebaseID,
		Notification: &messaging.Notification{
			Title: message.Title,
			Body:  message.Body,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
	})
	return err
}
