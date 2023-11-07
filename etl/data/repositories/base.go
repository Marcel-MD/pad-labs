package repositories

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T any] interface {
	FindUpdatedAfter(time.Time) ([]T, error)
	Upsert(t []T) error
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{
		db: db,
	}
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func (r *baseRepository[T]) FindUpdatedAfter(t time.Time) ([]T, error) {
	var ts []T
	err := r.db.Where("updated_at > ?", t).Find(&ts).Error
	return ts, err
}

func (r *baseRepository[T]) Upsert(t []T) error {
	if len(t) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&t).Error
}

func (r *baseRepository[T]) Delete(t *T) error {
	return r.db.Delete(t).Error
}
