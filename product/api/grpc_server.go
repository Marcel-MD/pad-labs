package api

import (
	"context"
	"fmt"
	"net"
	"product/api/mq"
	"product/api/pb"
	"product/config"
	"product/models"
	"product/services"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer(cfg config.Config, productService services.ProductService, producer mq.Producer, logger zerolog.Logger) (*grpc.Server, net.Listener, error) {
	log.Info().Msg("Creating new GRPC server")

	server := &grpcServer{
		productService: productService,
		producer:       producer,
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
	pb.RegisterProductServiceServer(srv, server)

	return srv, listener, nil
}

type grpcServer struct {
	pb.UnsafeProductServiceServer
	productService services.ProductService
	producer       mq.Producer
}

func mapProduct(product models.Product) *pb.Product {
	return &pb.Product{
		Id:        product.ID,
		OwnerId:   product.OwnerId,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: timestamppb.New(product.CreatedAt),
		UpdatedAt: timestamppb.New(product.UpdatedAt),
	}
}

func (s *grpcServer) GetAll(ctx context.Context, pbProductsQuery *pb.ProductsQuery) (*pb.Products, error) {
	productsQuery := models.PaginationQuery{
		Page: int(pbProductsQuery.Page),
		Size: int(pbProductsQuery.Size),
	}

	products, err := s.productService.FindAll(productsQuery)
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, mapProduct(product))
	}

	return &pb.Products{Products: pbProducts}, nil
}

func (s *grpcServer) GetById(ctx context.Context, pbProductId *pb.ProductId) (*pb.Product, error) {
	product, err := s.productService.FindById(pbProductId.Id)
	if err != nil {
		return nil, err
	}
	pbProduct := mapProduct(product)

	return pbProduct, nil
}

func (s *grpcServer) Create(ctx context.Context, pbCreateProduct *pb.CreateProduct) (*pb.ProductId, error) {
	product := models.Product{
		OwnerId: pbCreateProduct.OwnerId,
		Name:    pbCreateProduct.Name,
		Price:   pbCreateProduct.Price,
		Stock:   pbCreateProduct.Stock,
	}

	product, err := s.productService.Create(product)
	if err != nil {
		return nil, err
	}

	productMsg := models.ProductMessage{
		Base:    product.Base,
		OwnerId: product.OwnerId,
	}

	s.producer.SendMsg(models.CreateProductMsgType, productMsg, []string{models.OrderQueue})

	pbProductId := &pb.ProductId{Id: product.ID}

	return pbProductId, nil
}

func (s *grpcServer) Update(ctx context.Context, pbUpdateProduct *pb.UpdateProduct) (*emptypb.Empty, error) {
	product := models.Product{
		Name:    pbUpdateProduct.Name,
		Price:   pbUpdateProduct.Price,
		Stock:   pbUpdateProduct.Stock,
		OwnerId: pbUpdateProduct.OwnerId,
	}

	product.ID = pbUpdateProduct.Id

	product, err := s.productService.Update(product)
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
