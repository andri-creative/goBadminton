package midtrans

import (
	"os"
)

type Config struct {
	ServerKey string
	ClientKey string
	Env       string
}

func NewConfig() *Config {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	env := os.Getenv("MIDTRANS_ENV")
	if env == "" {
		env = "sandbox"
	}

	return &Config{
		ServerKey: serverKey,
		ClientKey: clientKey,
		Env:       env,
	}
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
