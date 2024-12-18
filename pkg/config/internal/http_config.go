package internal

type HTTP struct {
	Port string `envconfig:"PORT" default:"8000"`
}
