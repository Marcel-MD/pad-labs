package services

import (
	"errors"
	"order/config"
	"order/data/repositories"
	"order/models"

	"github.com/rs/zerolog/log"
)

type OrderService interface {
	FindAll(query models.PaginationQuery) ([]models.Order, error)
	FindById(id string) (models.Order, error)
	Create(order models.Order) (models.Order, error)
	Update(updateOrder models.UpdateOrder) (models.Order, error)
}

func NewOrderService(orderRepository repositories.OrderRepository, userRepository repositories.UserRepository, productRepository repositories.ProductRepository, cfg config.Config) OrderService {
	log.Info().Msg("Creating new order service")

	return &orderService{
		orderRepository:   orderRepository,
		userRepository:    userRepository,
		productRepository: productRepository,
		cfg:               cfg,
	}
}

type orderService struct {
	orderRepository   repositories.OrderRepository
	userRepository    repositories.UserRepository
	productRepository repositories.ProductRepository
	cfg               config.Config
}

func (s *orderService) FindAll(query models.PaginationQuery) ([]models.Order, error) {
	return s.orderRepository.FindAll(query)
}

func (s *orderService) FindById(id string) (models.Order, error) {
	return s.orderRepository.FindById(id)
}

func (s *orderService) Create(order models.Order) (models.Order, error) {
	_, err := s.userRepository.FindById(order.UserId)
	if err != nil {
		return order, err
	}

	_, err = s.productRepository.FindById(order.ProductId)
	if err != nil {
		return order, err
	}

	order.Status = models.PendingStatus

	err = s.orderRepository.Create(&order)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (s *orderService) Update(updateOrder models.UpdateOrder) (models.Order, error) {
	order, err := s.orderRepository.FindById(updateOrder.ID)
	if err != nil {
		return order, err
	}

	if order.Product.OwnerId != updateOrder.ProductOwnerId {
		return order, errors.New("not allowed to update this order")
	}

	order.Cost = updateOrder.Cost
	order.Status = updateOrder.Status

	err = s.orderRepository.Update(&order)
	if err != nil {
		return order, err
	}

	return order, nil
}
