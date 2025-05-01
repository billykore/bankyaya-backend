package internal

type App struct {
	Port string `envconfig:"PORT" default:"8000"`
}
