package main

import (
	"etl/config"
	"etl/data"
	"etl/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Warehouse
	warehouse, err := data.NewWarehouse(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to warehouse database")
	}

	// Source DBs
	userDB, err := data.NewUserDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to user database")
	}

	productDB, err := data.NewProductDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to product database")
	}

	orderDB, err := data.NewOrderDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to order database")
	}

	// Start Cron Jobs
	if err := services.StartCrons(cfg, warehouse, userDB, productDB, orderDB); err != nil {
		log.Fatal().Err(err).Msg("Failed to start cron jobs")
	}

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
	log.Warn().Msg("Shutting down ETL...")

	// Close DB connections
	if err := data.CloseDB(warehouse); err != nil {
		log.Err(err).Msg("Failed to close warehouse connection")
	}

	if err := data.CloseDB(userDB); err != nil {
		log.Err(err).Msg("Failed to close user database connection")
	}

	if err := data.CloseDB(productDB); err != nil {
		log.Err(err).Msg("Failed to close product database connection")
	}

	if err := data.CloseDB(orderDB); err != nil {
		log.Err(err).Msg("Failed to close order database connection")
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	responseJSON := []byte(`{"status": "ok"}`)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
