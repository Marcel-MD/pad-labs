package services

import (
	"errors"
	"testing"
	"time"

	"user/config"
	"user/mocks/mock_repositories"
	"user/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	registerUser := models.RegisterUser{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().FindByEmail(registerUser.Email).Return(models.User{}, errors.New("error"))
	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	token, err := userService.Register(registerUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	registerUser := models.RegisterUser{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().FindByEmail(registerUser.Email).Return(models.User{}, nil)

	token, err := userService.Register(registerUser)

	assert.EqualError(t, err, "user already exists")
	assert.Empty(t, token)
}

func TestUserService_Register_CreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	registerUser := models.RegisterUser{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().FindByEmail(registerUser.Email).Return(models.User{}, errors.New("error"))
	mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("error"))

	token, err := userService.Register(registerUser)

	assert.EqualError(t, err, "error")
	assert.Empty(t, token)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	loginUser := models.LoginUser{
		Email:    "test@example.com",
		Password: "password",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginUser.Password), bcrypt.DefaultCost)

	mockRepo.EXPECT().FindByEmail(loginUser.Email).Return(models.User{
		Base: models.Base{
			ID: "1",
		},
		Email:    loginUser.Email,
		Password: string(hashedPassword),
		Roles:    []string{models.UserRole},
	}, nil)

	token, err := userService.Login(loginUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, token.Token)
	assert.Equal(t, loginUser.Email, token.User.Email)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	loginUser := models.LoginUser{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().FindByEmail(loginUser.Email).Return(models.User{}, errors.New("user not found"))

	token, err := userService.Login(loginUser)

	assert.EqualError(t, err, "user not found")
	assert.Empty(t, token.Token)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockUserRepository(ctrl)
	cfg := config.Config{
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	userService := NewUserService(mockRepo, cfg)

	loginUser := models.LoginUser{
		Email:    "test@example.com",
		Password: "password",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("different_password"), bcrypt.DefaultCost)

	mockRepo.EXPECT().FindByEmail(loginUser.Email).Return(models.User{
		Base: models.Base{
			ID: "1",
		},
		Email:    loginUser.Email,
		Password: string(hash),
		Roles:    []string{models.UserRole},
	}, nil)

	token, err := userService.Login(loginUser)

	assert.EqualError(t, err, "crypto/bcrypt: hashedPassword is not the hash of the given password")
	assert.Empty(t, token.Token)
}
