package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port      string
	ClientUrl string
	LogLevel  string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type PrometheusConfig struct {
	Addr string
}

type Config struct {
	Server     ServerConfig
	Redis      RedisConfig
	Prometheus PrometheusConfig
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("server.port", ":8080")
	v.SetDefault("server.client_url", "http://localhost:3000")
	v.SetDefault("server.log_level", "info")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	v.SetEnvPrefix("CONT-DSH")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("no config file is there")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:      v.GetString("server.port"),
			ClientUrl: v.GetString("server.client_url"),
			LogLevel:  v.GetString("server.log_level"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("redis.addr"),
			Password: v.GetString("redis.password"),
		},
		Prometheus: PrometheusConfig{
			Addr: v.GetString("prometheus.addr"),
		},
	}

	return cfg, nil

}
