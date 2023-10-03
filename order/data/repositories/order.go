package repositories

import (
	"order/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll(query models.PaginationQuery) ([]models.Order, error)
	FindById(id string) (models.Order, error)
	Create(t *models.Order) error
	Update(t *models.Order) error
	Delete(t *models.Order) error
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	log.Info().Msg("Creating new order repository")

	return &orderRepository{
		BaseRepository: NewBaseRepository[models.Order](db),
		db:             db,
	}
}

type orderRepository struct {
	BaseRepository[models.Order]
	db *gorm.DB
}

func (r *orderRepository) FindById(id string) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("Product").First(&order, "id = ?", id).Error
	if err != nil {
		return order, err
	}

	return order, nil
}
