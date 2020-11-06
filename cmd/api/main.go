package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mat-greaves/test-api/internal/config"
	"github.com/mat-greaves/test-api/internal/models"
	"github.com/mat-greaves/test-api/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cfg := config.GetConfig()
	validate := validator.New()
	db, err := models.NewDB("mongodb://localhost:27017")
	if err != nil {
		logger.Fatal().Msgf("failed to connect to database: %s", err.Error())
	}
	logger.Info().Msg("connected to database")
	s := server.NewServer(&logger, cfg, db, validate)

	s.Run()
}
