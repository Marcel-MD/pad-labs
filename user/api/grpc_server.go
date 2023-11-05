package api

import (
	"context"
	"fmt"
	"net"
	"user/api/mq"
	"user/api/pb"
	"user/config"
	"user/models"
	"user/services"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer(cfg config.Config, userService services.UserService, producer mq.Producer, logger zerolog.Logger) (*grpc.Server, net.Listener, *prometheus.Registry, error) {
	log.Info().Msg("Creating new GRPC server")

	// Setup metrics.
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	server := &grpcServer{
		userService: userService,
		producer:    producer,
	}

	listener, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		return nil, nil, nil, err
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)),
			srvMetrics.UnaryServerInterceptor(),
		),
	)
	pb.RegisterUserServiceServer(srv, server)

	return srv, listener, reg, nil
}

type grpcServer struct {
	pb.UnsafeUserServiceServer
	userService services.UserService
	producer    mq.Producer
}

func mapUser(user models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
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

	s.producer.SendMsg(models.CreateUserMsgType, token.User.Base, []string{models.ProductQueue, models.OrderQueue})

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

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
