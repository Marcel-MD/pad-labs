package repositories

import (
	"product/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Create(t *models.User) error
	Update(t *models.User) error
	Delete(t *models.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	log.Info().Msg("Creating new user repository")

	return &userRepository{
		BaseRepository: NewBaseRepository[models.User](db),
		db:             db,
	}
}

type userRepository struct {
	BaseRepository[models.User]
	db *gorm.DB
}
