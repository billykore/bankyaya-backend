package internal

// Email config.
type Email struct {
	From     string `envconfig:"EMAIL_FROM"`
	Host     string `envconfig:"EMAIL_HOST"`
	Port     int    `envconfig:"EMAIL_PORT"`
	Username string `envconfig:"EMAIL_USERNAME"`
	Password string `envconfig:"EMAIL_PASSWORD"`
}
