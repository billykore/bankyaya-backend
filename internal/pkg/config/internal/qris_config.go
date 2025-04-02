package internal

type QRIS struct {
	URL          string `envconfig:"QRIS_URL"`
	ClientId     string `envconfig:"QRIS_CLIENT_ID"`
	ClientSecret string `envconfig:"QRIS_CLIENT_SECRET"`
}
