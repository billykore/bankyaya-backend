package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
)

type Connection struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewConnection establishes a new RabbitMQ connection.
func NewConnection(cfg *config.Config) *Connection {
	log := logger.New()

	conn, err := amqp.Dial(cfg.Rabbit.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a Channel: %s", err)
		return nil
	}
	return &Connection{
		conn:    conn,
		Channel: channel,
	}
}

// Close closes the RabbitMQ connection and channel.
func (c *Connection) Close() {
	err := c.Channel.Close()
	if err != nil {
		return
	}
	err = c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Connection) Publish(ctx context.Context, queue string, body []byte) error {
	err := c.Channel.ExchangeDeclare(
		queue,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	// Publish message to the RabbitMQ exchange
	err = c.Channel.PublishWithContext(ctx,
		queue, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Consume(ctx context.Context, queue string) (<-chan amqp.Delivery, error) {
	err := c.Channel.ExchangeDeclare(
		queue,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := c.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = c.Channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		queue,  // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := c.Channel.ConsumeWithContext(ctx,
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
