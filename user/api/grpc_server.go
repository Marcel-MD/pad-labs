package api

import (
	"context"
	"net"
	"user/api/pb"
	"user/config"
	"user/models"
	"user/services"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer(cfg config.Config, userService services.UserService) (*grpc.Server, net.Listener, error) {
	log.Info().Msg("Creating new GRPC server")

	server := &grpcServer{
		userService: userService,
	}

	listener, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, server)

	return srv, listener, nil
}

type grpcServer struct {
	pb.UnsafeUserServiceServer
	userService services.UserService
}

func mapUser(user models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func (s *grpcServer) Register(ctx context.Context, pbRegisterUser *pb.RegisterUser) (*pb.Token, error) {
	registerUser := models.RegisterUser{
		Email:    pbRegisterUser.Email,
		Name:     pbRegisterUser.Name,
		Password: pbRegisterUser.Password,
	}

	token, err := s.userService.Register(registerUser)
	if err != nil {
		return nil, err
	}

	pbToken := pb.Token{
		Token: token.Token,
	}

	return &pbToken, nil
}

func (s *grpcServer) Login(ctx context.Context, pbLoginUser *pb.LoginUser) (*pb.Token, error) {
	loginUser := models.LoginUser{
		Email:    pbLoginUser.Email,
		Password: pbLoginUser.Password,
	}

	token, err := s.userService.Login(loginUser)
	if err != nil {
		return nil, err
	}

	pbToken := pb.Token{
		Token: token.Token,
	}

	return &pbToken, nil
}

func (s *grpcServer) Validate(ctx context.Context, pbToken *pb.Token) (*pb.User, error) {
	user, err := s.userService.Validate(pbToken.Token)
	if err != nil {
		return nil, err
	}
	pbUser := mapUser(user)

	return pbUser, nil
}

func (s *grpcServer) GetAll(ctx context.Context, pbUsersQuery *pb.UsersQuery) (*pb.Users, error) {
	usersQuery := models.PaginationQuery{
		Page: int(pbUsersQuery.Page),
		Size: int(pbUsersQuery.Size),
	}

	users, err := s.userService.FindAll(usersQuery)
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, mapUser(user))
	}

	return &pb.Users{Users: pbUsers}, nil
}

func (s *grpcServer) GetById(ctx context.Context, pbUserId *pb.UserId) (*pb.User, error) {
	user, err := s.userService.FindById(pbUserId.Id)
	if err != nil {
		return nil, err
	}
	pbUser := mapUser(user)

	return pbUser, nil
}

func (s *grpcServer) Delete(ctx context.Context, pbUserId *pb.UserId) (*emptypb.Empty, error) {
	err := s.userService.Delete(pbUserId.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
