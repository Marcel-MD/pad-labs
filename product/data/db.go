package data

import (
	"product/config"
	"product/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.Config) (*gorm.DB, error) {
	log.Info().Msg("Creating new database connection")

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	dbSql, err := db.DB()
	if err != nil {
		return err
	}

	return dbSql.Close()
}
