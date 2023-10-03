package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"user/api"
	"user/api/mq"
	"user/config"
	"user/data"
	"user/data/repositories"
	"user/services"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// DB
	db, err := data.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// RabbitMQ
	producer, err := mq.NewProducer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create RabbitMQ producer")
	}

	// User
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, cfg)

	// Start GRPC Server
	grpcSrv, listener, err := api.NewGrpcServer(cfg, userService, producer)
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
	go func() {
		if err := http.ListenAndServe(cfg.HttpPort, nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
	log.Warn().Msg("Shutting down user server...")

	// Shutdown GRPC server
	grpcSrv.GracefulStop()

	// Close RabbitMQ connection
	if err := producer.Close(); err != nil {
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
