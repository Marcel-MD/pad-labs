package services

import (
	"order/config"
	"order/data/repositories"
	"order/models"

	"github.com/rs/zerolog/log"
)

type ProductService interface {
	FindAll(query models.PaginationQuery) ([]models.Product, error)
	FindById(id string) (models.Product, error)
	Create(product models.Product) (models.Product, error)
}

func NewProductService(productRepository repositories.ProductRepository, userRepository repositories.UserRepository, cfg config.Config) ProductService {
	log.Info().Msg("Creating new product service")

	return &productService{
		productRepository: productRepository,
		userRepository:    userRepository,
		cfg:               cfg,
	}
}

type productService struct {
	productRepository repositories.ProductRepository
	userRepository    repositories.UserRepository
	cfg               config.Config
}

func (s *productService) FindAll(query models.PaginationQuery) ([]models.Product, error) {
	return s.productRepository.FindAll(query)
}

func (s *productService) FindById(id string) (models.Product, error) {
	return s.productRepository.FindById(id)
}

func (s *productService) Create(product models.Product) (models.Product, error) {
	_, err := s.userRepository.FindById(product.OwnerId)
	if err != nil {
		return product, err
	}

	err = s.productRepository.Create(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}
