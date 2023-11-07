package services

import (
	"etl/config"
	"etl/data/repositories"
	"etl/models"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func StartCrons(cfg config.Config, whDB, userDB, productDB, orderDB *gorm.DB) error {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(cfg.SyncInterval).Do(func() {
		if err := syncDB[models.User](models.UserDB, userDB, whDB); err != nil {
			log.Err(err).Msg("Failed to sync user database")
		}

		if err := syncDB[models.Product](models.ProductDB, productDB, whDB); err != nil {
			log.Err(err).Msg("Failed to sync product database")
		}

		if err := syncDB[models.Order](models.OrderDB, orderDB, whDB); err != nil {
			log.Err(err).Msg("Failed to sync order database")
		}
	})

	if err != nil {
		return err
	}

	s.StartAsync()
	return nil
}

func syncDB[T any](databaseName string, sourceDB, whDB *gorm.DB) error {
	log.Info().Str("database", databaseName).Msg("Syncing database")

	stateRepo := repositories.NewStateRepository(whDB)
	sourceRepo := repositories.NewBaseRepository[T](sourceDB)
	whRepo := repositories.NewBaseRepository[T](whDB)

	lastSync, err := stateRepo.GetLastSync(databaseName)
	if err != nil {
		log.Err(err).Msg("Failed to get last sync timestamp")
		lastSync = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	data, err := sourceRepo.FindUpdatedAfter(lastSync)
	if err != nil {
		return err
	}

	err = whRepo.Upsert(data)
	if err != nil {
		return err
	}

	err = stateRepo.SetLastSync(databaseName, time.Now())
	if err != nil {
		log.Err(err).Msg("Failed to set last sync timestamp")
	}

	return nil
}
