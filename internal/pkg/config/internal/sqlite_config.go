package internal

type SQLite struct {
	Name string `envconfig:"SQLITE_DB_NAME"`
}
