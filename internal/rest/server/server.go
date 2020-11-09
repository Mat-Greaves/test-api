// HTTP Server definition
package server

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mat-greaves/test-api/internal/config"
	"github.com/mat-greaves/test-api/internal/rest/handlers"
	"github.com/mat-greaves/test-api/internal/rest/middleware"
	"github.com/mat-greaves/test-api/internal/service"
	"github.com/rs/zerolog"
)

type Server struct {
	logger   *zerolog.Logger
	cfg      config.Configer
	validate *validator.Validate
	service  *service.Service
}

func NewServer(
	logger *zerolog.Logger,
	cfg config.Configer,
	validate *validator.Validate,
	service *service.Service,
) *Server {
	return &Server{
		logger:   logger,
		cfg:      cfg,
		validate: validate,
		service:  service,
	}
}

func (s *Server) Run() {
	// handler has http request handlers
	userHandler := handlers.NewUserHandler(s.validate, s.service.Users)
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
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
