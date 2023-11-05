package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"order/api"
	"order/api/mq"
	"order/config"
	"order/data"
	"order/data/repositories"
	"order/services"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Logger
	logger := zerolog.New(os.Stderr)

	// DB
	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// RabbitMQ Producer
	producer, err := mq.NewProducer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create RabbitMQ producer")
	}

	// User
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, cfg)

	// Product
	productRepository := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepository, userRepository, cfg)

	// Order
	orderRepository := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepository, userRepository, productRepository, cfg)

	// RabbitMQ Consumer
	consumer, err := mq.NewConsumer(cfg, userService, productService, orderService)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create RabbitMQ consumer")
	}

	// Start GRPC Server
	grpcSrv, listener, reg, err := api.NewGrpcServer(cfg, orderService, producer, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create gRPC server")
	}
	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("Failed to start gRPC server")
		}
	}()

	// Start HTTP Server
	http.HandleFunc("/", statusHandler)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		if err := http.ListenAndServe(cfg.HttpPort, nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
	log.Warn().Msg("Shutting down product server...")

	// Shutdown GRPC server
	grpcSrv.GracefulStop()

	// Close RabbitMQ producer
	if err := producer.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close RabbitMQ connection")
	}

	// Close RabbitMQ consumer
	if err := consumer.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close RabbitMQ connection")
	}

	// Close DB connection
	if err := data.CloseDB(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to close db connection")
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	responseJSON := []byte(`{"status": "ok"}`)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
