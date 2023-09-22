package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user/api"
	"user/api/controllers"
	"user/api/mq"
	"user/config"
	"user/data"
	"user/data/repositories"
	"user/services"

	"github.com/rs/zerolog/log"
)

// @title User API
// @description This is user server.
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
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
	userController := controllers.NewUserController(userService, producer)

	// Start HTTP Server
	httpSrv := api.NewHttpServer(cfg, userController)
	go func() {
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
		log.Info().Msg("All HTTP server connections are closed")
	}()

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
	log.Warn().Msg("Shutting down HTTP server...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("HTTP Server forced to shutdown")
	}

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
