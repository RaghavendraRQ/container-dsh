package config

import (
	"fmt"
	"os"
)

type Cfg struct {
	PORT       string
	CLIENT_URL string
}

type redisCfg struct {
	Addr     string
	Password string
}

func NewConfig() (*Cfg, error) {
	PORT, ok := os.LookupEnv("PORT")
	if !ok {
		PORT = ":8080"
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

func NewRedisConfig() (*redisCfg, error) {
	rdb := &redisCfg{}
	Addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		rdb.Addr = "localhost:6379"
	} else {
		rdb.Addr = Addr
	}

	Password, ok := os.LookupEnv("REDIS_PASSWD")
	if !ok {
		rdb.Password = ""
	} else {
		rdb.Password = Password
	}

	return rdb, nil

}
