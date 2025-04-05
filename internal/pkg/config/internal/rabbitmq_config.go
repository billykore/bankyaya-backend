package internal

// Rabbit configuration.
type Rabbit struct {
	DSN       string `envconfig:"RABBIT_DSN"`
	QueueName string `envconfig:"RABBIT_QUEUE_NAME"`
}
