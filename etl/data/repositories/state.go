package repositories

import (
	"etl/models"
	"time"

	"gorm.io/gorm"
)

type StateRepository interface {
	GetLastSync(databaseName string) (time.Time, error)
	SetLastSync(databaseName string, t time.Time) error
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return &stateRepository{
		db: db,
	}
}

type stateRepository struct {
	db *gorm.DB
}

func (r *stateRepository) GetLastSync(databaseName string) (time.Time, error) {
	var state models.State
	err := r.db.First(&state, "database_name = ?", databaseName).Error
	if err != nil {
		return time.Time{}, err
	}

	return state.LastSync, nil
}

func (r *stateRepository) SetLastSync(databaseName string, t time.Time) error {
	var state models.State
	err := r.db.Where(models.State{DatabaseName: databaseName}).FirstOrCreate(&state).Error
	if err != nil {
		return err
	}

	state.LastSync = t
	return r.db.Save(&state).Error
}
