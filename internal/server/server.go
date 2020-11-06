package server

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mat-greaves/test-api/internal/config"
	"github.com/mat-greaves/test-api/internal/handlers"
	"github.com/mat-greaves/test-api/internal/middleware"
	"github.com/mat-greaves/test-api/internal/models"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	logger   *zerolog.Logger
	cfg      config.Configer
	db       *mongo.Client
	validate *validator.Validate
}

func NewServer(
	logger *zerolog.Logger,
	cfg config.Configer,
	db *mongo.Client,
	validate *validator.Validate,
) *Server {
	return &Server{
		logger:   logger,
		cfg:      cfg,
		db:       db,
		validate: validate,
	}
}

func (s *Server) Run() {
	database := s.db.Database("test-api")
	// store has methods for interacting with users in database
	userStore := models.NewUserStore(database.Collection("users"))
	// handler has http request handlers
	userHandler := handlers.NewUserHandler(s.validate, userStore)
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	// mux middleware only run on matches, we want requestId and logging middleware to also run on NotFound
	// in order to achieve this we wrap the router to execute the middleware before it even looks for a match
	// https://github.com/gorilla/mux/issues/416
	requestIdMiddleware := middleware.UseRequestId()
	loggingMiddleware := middleware.UseLogger(s.logger)
	loggedRouter := requestIdMiddleware(loggingMiddleware(r))
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		s.logger.Fatal().Msgf("Error listenAndServe fail: %s\n", err)
	}
}
