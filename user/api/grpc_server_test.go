package api

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"user/api/pb"
	"user/config"
	"user/mocks/mock_mq"
	"user/mocks/mock_services"
	"user/models"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

func genRandomPort() string {
	num := 50000 + rand.Intn(1000)
	return fmt.Sprintf(":%d", num)
}

func TestGrpcServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockProducer := mock_mq.NewMockProducer(ctrl)
	cfg := config.Config{
		GrpcPort:            genRandomPort(),
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	server, listener, err := NewGrpcServer(cfg, mockUserService, mockProducer, zerolog.New(os.Stderr))
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(cfg.GrpcPort, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	registerUser := models.RegisterUser{
		Email:    "test@example.com",
		Password: "password",
	}

	token := models.Token{
		Token: "token",
		User: models.User{
			Base: models.Base{
				ID:        "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "test@example.com",
			Name:  "Test User",
			Roles: []string{"user"},
		},
	}

	mockUserService.EXPECT().Register(registerUser).Return(token, nil)
	mockProducer.EXPECT().SendMsg(models.CreateUserMsgType, token.User.Base, []string{models.ProductQueue, models.OrderQueue})

	pbRegisterUser := &pb.RegisterUser{
		Email:    registerUser.Email,
		Name:     registerUser.Name,
		Password: registerUser.Password,
	}

	pbToken, err := client.Register(context.Background(), pbRegisterUser)
	require.NoError(t, err)

	expectedPbToken := &pb.Token{
		Token: token.Token,
	}

	assert.Equal(t, expectedPbToken.Token, pbToken.Token)
}

func TestGrpcServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockProducer := mock_mq.NewMockProducer(ctrl)
	cfg := config.Config{
		GrpcPort:            genRandomPort(),
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	server, listener, err := NewGrpcServer(cfg, mockUserService, mockProducer, zerolog.New(os.Stderr))
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(cfg.GrpcPort, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	loginUser := models.LoginUser{
		Email:    "test@example.com",
		Password: "password",
	}

	token := models.Token{
		Token: "token",
		User: models.User{
			Base: models.Base{
				ID:        "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "test@example.com",
			Name:  "Test User",
			Roles: []string{"user"},
		},
	}

	mockUserService.EXPECT().Login(loginUser).Return(token, nil)

	pbLoginUser := &pb.LoginUser{
		Email:    loginUser.Email,
		Password: loginUser.Password,
	}

	pbToken, err := client.Login(context.Background(), pbLoginUser)
	require.NoError(t, err)

	expectedPbToken := &pb.Token{
		Token: token.Token,
	}

	assert.Equal(t, expectedPbToken.Token, pbToken.Token)
}

func TestGrpcServer_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockProducer := mock_mq.NewMockProducer(ctrl)
	cfg := config.Config{
		GrpcPort:            genRandomPort(),
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	server, listener, err := NewGrpcServer(cfg, mockUserService, mockProducer, zerolog.New(os.Stderr))
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(cfg.GrpcPort, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	user := models.User{
		Base: models.Base{
			ID:        "1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "test@example.com",
		Name:  "Test User",
		Roles: []string{"user"},
	}

	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}

	token := models.Token{
		Token: "token",
		User:  user,
	}

	mockUserService.EXPECT().Validate(token.Token).Return(user, nil)

	pbToken := &pb.Token{
		Token: token.Token,
	}

	pbUser, err = client.Validate(context.Background(), pbToken)
	require.NoError(t, err)

	expectedPbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}

	assert.Equal(t, expectedPbUser.Id, pbUser.Id)
	assert.Equal(t, expectedPbUser.Email, pbUser.Email)
	assert.Equal(t, expectedPbUser.Name, pbUser.Name)
}

func TestGrpcServer_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockProducer := mock_mq.NewMockProducer(ctrl)
	cfg := config.Config{
		GrpcPort:            genRandomPort(),
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	server, listener, err := NewGrpcServer(cfg, mockUserService, mockProducer, zerolog.New(os.Stderr))
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(cfg.GrpcPort, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	users := []models.User{
		{
			Base: models.Base{
				ID:        "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "test1@example.com",
			Name:  "Test User 1",
			Roles: []string{"user"},
		},
		{
			Base: models.Base{
				ID:        "2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "test2@example.com",
			Name:  "Test User 2",
			Roles: []string{"user"},
		},
	}

	pbUsers := []*pb.User{
		{
			Id:        users[0].ID,
			Email:     users[0].Email,
			Name:      users[0].Name,
			Roles:     users[0].Roles,
			CreatedAt: &timestamp.Timestamp{Seconds: users[0].CreatedAt.Unix()},
			UpdatedAt: &timestamp.Timestamp{Seconds: users[0].UpdatedAt.Unix()},
		},
		{
			Id:        users[1].ID,
			Email:     users[1].Email,
			Name:      users[1].Name,
			Roles:     users[1].Roles,
			CreatedAt: &timestamp.Timestamp{Seconds: users[1].CreatedAt.Unix()},
			UpdatedAt: &timestamp.Timestamp{Seconds: users[1].UpdatedAt.Unix()},
		},
	}

	usersQuery := models.PaginationQuery{
		Page: 1,
		Size: 10,
	}

	mockUserService.EXPECT().FindAll(usersQuery).Return(users, nil)

	pbUsersQuery := &pb.UsersQuery{
		Page: int64(usersQuery.Page),
		Size: int64(usersQuery.Size),
	}

	pbUsersResponse, err := client.GetAll(context.Background(), pbUsersQuery)
	require.NoError(t, err)

	expectedPbUsers := &pb.Users{
		Users: pbUsers,
	}

	for i := range expectedPbUsers.Users {
		assert.Equal(t, expectedPbUsers.Users[i].Id, pbUsersResponse.Users[i].Id)
		assert.Equal(t, expectedPbUsers.Users[i].Email, pbUsersResponse.Users[i].Email)
		assert.Equal(t, expectedPbUsers.Users[i].Name, pbUsersResponse.Users[i].Name)
	}
}

func TestGrpcServer_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockProducer := mock_mq.NewMockProducer(ctrl)
	cfg := config.Config{
		GrpcPort:            genRandomPort(),
		AccessTokenSecret:   "SecretSecretSecret",
		AccessTokenLifespan: 72 * time.Hour,
	}

	server, listener, err := NewGrpcServer(cfg, mockUserService, mockProducer, zerolog.New(os.Stderr))
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(cfg.GrpcPort, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	user := models.User{
		Base: models.Base{
			ID:        "1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "test@example.com",
		Name:  "Test User",
		Roles: []string{"user"},
	}

	pbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}

	mockUserService.EXPECT().FindById(user.ID).Return(user, nil)

	pbUserId := &pb.UserId{
		Id: user.ID,
	}

	pbUser, err = client.GetById(context.Background(), pbUserId)
	require.NoError(t, err)

	expectedPbUser := &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}

	assert.Equal(t, expectedPbUser.Id, pbUser.Id)
	assert.Equal(t, expectedPbUser.Email, pbUser.Email)
	assert.Equal(t, expectedPbUser.Name, pbUser.Name)
}
