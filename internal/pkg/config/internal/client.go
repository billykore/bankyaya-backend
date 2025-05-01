package internal

type Clients []Client

type Client struct {
	Name       string
	MinVersion string
	MaxVersion string
}
