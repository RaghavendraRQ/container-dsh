package config

import (
	"fmt"
	"os"
)

type Cfg struct {
	PORT       string
	CLIENT_URL string
}

func NewConfig() (*Cfg, error) {
	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = "8080"
	}

	CLIENT_URL, ok := os.LookupEnv("CLIENT_URL")
	if !ok {
		return nil, fmt.Errorf("client url is unset")
	}

	return &Cfg{
		PORT,
		CLIENT_URL,
	}, nil

}
