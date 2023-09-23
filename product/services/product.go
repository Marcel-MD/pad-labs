package services

import (
	"product/config"
	"product/data/repositories"
	"product/models"

	"github.com/rs/zerolog/log"
)

type ProductService interface {
	FindAll(query models.PaginationQuery) ([]models.Product, error)
	FindById(id string) (models.Product, error)
	Create(product models.Product) (models.Product, error)
	Delete(id string) error
}

func NewProductService(repository repositories.ProductRepository, cfg config.Config) ProductService {
	log.Info().Msg("Creating new product service")

	return &productService{
		repository: repository,
		cfg:        cfg,
	}
}

type productService struct {
	repository repositories.ProductRepository
	cfg        config.Config
}

func (s *productService) FindAll(query models.PaginationQuery) ([]models.Product, error) {
	return s.repository.FindAll(query)
}

func (s *productService) FindById(id string) (models.Product, error) {
	return s.repository.FindById(id)
}

func (s *productService) Delete(id string) error {
	product, err := s.repository.FindById(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(&product)
}

func (s *productService) Create(product models.Product) (models.Product, error) {
	err := s.repository.Create(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}
