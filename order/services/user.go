package services

import (
	"order/config"
	"order/data/repositories"
	"order/models"

	"github.com/rs/zerolog/log"
)

type UserService interface {
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Create(user models.User) (models.User, error)
}

func NewUserService(repository repositories.UserRepository, cfg config.Config) UserService {
	log.Info().Msg("Creating new user service")

	return &userService{
		repository: repository,
		cfg:        cfg,
	}
}

type userService struct {
	repository repositories.UserRepository
	cfg        config.Config
}

func (s *userService) FindAll(query models.PaginationQuery) ([]models.User, error) {
	return s.repository.FindAll(query)
}

func (s *userService) FindById(id string) (models.User, error) {
	return s.repository.FindById(id)
}

func (s *userService) Create(user models.User) (models.User, error) {
	err := s.repository.Create(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}
