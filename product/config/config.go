package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpPort    string `env:"HTTP_PORT" envDefault:":8090"`
	GrpcPort    string `env:"GRPC_PORT" envDefault:":8091"`
	Env         string `env:"ENV" envDefault:"dev"`
	RabbitMQUrl string `env:"RABBIT_MQ_URL" envDefault:"amqp://guest:guest@rabbitmq:5672/"`
	DatabaseUrl string `env:"DATABASE_URL" envDefault:"postgres://postgres:password@product-db:5432/product-db"`
}

func NewConfig() (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file.")
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
