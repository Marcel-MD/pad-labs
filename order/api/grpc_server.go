package api

import (
	"context"
	"fmt"
	"net"
	"order/api/mq"
	"order/api/pb"
	"order/config"
	"order/models"
	"order/services"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer(cfg config.Config, orderService services.OrderService, producer mq.Producer, logger zerolog.Logger) (*grpc.Server, net.Listener, error) {
	log.Info().Msg("Creating new GRPC server")

	server := &grpcServer{
		orderService: orderService,
		producer:     producer,
	}

	listener, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)),
		),
	)
	pb.RegisterOrderServiceServer(srv, server)

	return srv, listener, nil
}

type grpcServer struct {
	pb.UnsafeOrderServiceServer
	orderService services.OrderService
	producer     mq.Producer
}

func mapOrder(order models.Order) *pb.Order {
	return &pb.Order{
		Id:              order.ID,
		ProductId:       order.ProductId,
		UserId:          order.UserId,
		Quantity:        order.Quantity,
		Cost:            order.Cost,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}
}

func (s *grpcServer) GetAll(ctx context.Context, pbOrdersQuery *pb.OrdersQuery) (*pb.Orders, error) {
	ordersQuery := models.PaginationQuery{
		Page: int(pbOrdersQuery.Page),
		Size: int(pbOrdersQuery.Size),
	}

	orders, err := s.orderService.FindAll(ordersQuery)
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, mapOrder(order))
	}

	return &pb.Orders{Orders: pbOrders}, nil
}

func (s *grpcServer) GetById(ctx context.Context, pbOrderId *pb.OrderId) (*pb.Order, error) {
	order, err := s.orderService.FindById(pbOrderId.Id)
	if err != nil {
		return nil, err
	}
	pbOrder := mapOrder(order)

	return pbOrder, nil
}

func (s *grpcServer) Create(ctx context.Context, pbCreateOrder *pb.CreateOrder) (*pb.OrderId, error) {
	order := models.Order{
		ProductId:       pbCreateOrder.ProductId,
		UserId:          pbCreateOrder.UserId,
		Quantity:        pbCreateOrder.Quantity,
		ShippingAddress: pbCreateOrder.ShippingAddress,
	}

	order, err := s.orderService.Create(order)
	if err != nil {
		return nil, err
	}

	s.producer.SendMsg(models.CreateOrderMsgType, order, []string{models.ProductQueue})

	pbOrderId := &pb.OrderId{Id: order.ID}

	return pbOrderId, nil
}

func (s *grpcServer) Update(ctx context.Context, pbUpdateOrder *pb.UpdateOrder) (*emptypb.Empty, error) {
	updateOrder := models.UpdateOrder{
		ID:             pbUpdateOrder.Id,
		ProductOwnerId: pbUpdateOrder.ProductOwnerId,
		Status:         pbUpdateOrder.Status,
		Cost:           pbUpdateOrder.Cost,
	}

	_, err := s.orderService.Update(updateOrder)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
