package internal

type CoreBanking struct {
	URL      string `envconfig:"COREBANKING_URL"`
	Username string `envconfig:"COREBANKING_USERNAME"`
	Password string `envconfig:"COREBANKING_PASSWORD"`
}
