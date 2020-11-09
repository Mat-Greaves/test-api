package main

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mat-greaves/test-api/internal/config"
	"github.com/mat-greaves/test-api/internal/rest/server"
	"github.com/mat-greaves/test-api/internal/service"
	"github.com/mat-greaves/test-api/internal/service/users"
	"github.com/mat-greaves/test-api/internal/store"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cfg := config.GetConfig()
	validate := validator.New()
	db, err := store.NewDB("mongodb://localhost:27017")
	if err != nil {
		logger.Fatal().Msgf("failed to connect to database: %s", err.Error())
	}
	logger.Info().Msg("connected to database")
	database := db.Database("test-api")
	userStore := users.NewUserStore(database.Collection("users"))
	userService := users.NewUserService(userStore)
	svc := service.Service{
		Users: userService,
	}
	s := server.NewServer(&logger, cfg, validate, &svc)
	s.Run()
}
