package config

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env          string        `env:"ENV" envDefault:"dev"`
	HttpPort     string        `env:"HTTP_PORT" envDefault:":8060"`
	SyncInterval time.Duration `env:"SYNC_INTERVAL" envDefault:"5m"`

	UserDbUrl    string `env:"USER_DB_URL" envDefault:"postgres://postgres:password@user-db:5432/user-db"`
	ProductDbUrl string `env:"PRODUCT_DB_URL" envDefault:"postgres://postgres:password@product-db:5432/product-db"`
	OrderDbUrl   string `env:"ORDER_DB_URL" envDefault:"postgres://postgres:password@order-db:5432/order-db"`
	WarehouseUrl string `env:"WAREHOUSE_URL" envDefault:"postgres://postgres:password@warehouse:5432/warehouse"`
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
