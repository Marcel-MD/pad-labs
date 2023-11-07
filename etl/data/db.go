package data

import (
	"etl/config"
	"etl/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewWarehouse(cfg config.Config) (*gorm.DB, error) {
	log.Info().Msg("Creating new Warehouse connection")

	db, err := gorm.Open(postgres.Open(cfg.WarehouseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.State{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Order{})

	return db, nil
}

func NewUserDB(cfg config.Config) (*gorm.DB, error) {
	log.Info().Msg("Creating new User database connection")

	db, err := gorm.Open(postgres.Open(cfg.UserDbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewProductDB(cfg config.Config) (*gorm.DB, error) {
	log.Info().Msg("Creating new Product database connection")

	db, err := gorm.Open(postgres.Open(cfg.ProductDbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewOrderDB(cfg config.Config) (*gorm.DB, error) {
	log.Info().Msg("Creating new Order database connection")

	db, err := gorm.Open(postgres.Open(cfg.OrderDbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	dbSql, err := db.DB()
	if err != nil {
		return err
	}

	return dbSql.Close()
}
