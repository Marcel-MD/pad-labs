package api

import (
	"context"
	"net"
	"product/api/mq"
	"product/api/pb"
	"product/config"
	"product/models"
	"product/services"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGrpcServer(cfg config.Config, productService services.ProductService, producer mq.Producer) (*grpc.Server, net.Listener, error) {
	log.Info().Msg("Creating new GRPC server")

	server := &grpcServer{
		productService: productService,
		producer:       producer,
	}

	listener, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()
	pb.RegisterProductServiceServer(srv, server)

	return srv, listener, nil
}

type grpcServer struct {
	// pb.UnsafeProductServiceServer
	pb.ProductServiceServer
	productService services.ProductService
	producer       mq.Producer
}

func mapProduct(product models.Product) *pb.Product {
	return &pb.Product{
		Id:        product.ID,
		Name:      product.Name,
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
