package repositories

import (
	"order/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(query models.PaginationQuery) ([]models.Product, error)
	FindById(id string) (models.Product, error)
	Create(t *models.Product) error
	Update(t *models.Product) error
	Delete(t *models.Product) error
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	log.Info().Msg("Creating new product repository")

	return &productRepository{
		BaseRepository: NewBaseRepository[models.Product](db),
		db:             db,
	}
}

type productRepository struct {
	BaseRepository[models.Product]
	db *gorm.DB
}
