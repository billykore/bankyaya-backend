package internal

// Token config.
type Token struct {
	Secret    string `envconfig:"TOKEN_SECRET"`
	HeaderKid string `envconfig:"TOKEN_HEADER_KEY_ID"`
}
