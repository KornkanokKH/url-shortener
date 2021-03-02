package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/cmd/url-shortener/generate"
	"url-shortener/internal/config"
)

func main() {
	log.Info().Msg("start main...")

	if err := run(); err != nil {
		log.Error().Msgf("Unexpected error to run server: %v", err)
		os.Exit(1)
	}
}
func run() error {

	// Read Config
	conf := new(Config)
	if err := config.ReadConfigFile(conf); err != nil {
		log.Fatal().Msgf("Unexpected error to init configuration: %v.", err)
	}

	generate := generate.NewService()

	server := &http.Server{
		Handler:      routes(generate),
		Addr:         fmt.Sprintf(":%v", conf.Port),
		WriteTimeout: conf.Timeout * time.Second,
		ReadTimeout:  conf.Timeout * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()
	log.Info().Msg("Server Started")

	<-done
	log.Info().Msg("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Server Shutdown Failed:%+v", err)
	}
	log.Info().Msg("Server Exited Properly")

	return nil
}

func routes(generate generate.Service) *mux.Router {

	route := mux.NewRouter()

	route.HandleFunc("/generate", generate.GenerateUrlShortener).Methods(http.MethodPost)
	route.StrictSlash(false)

	return nrgorilla.InstrumentRoutes(route, nil)
}