package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func WriteJSON(l *zerolog.Logger, w http.ResponseWriter, status int, msg interface{}) bool {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if msg == nil {
		return true
	}
	data, err := json.Marshal(msg)
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Error marshaling response"))
		return false
	}
	_, err = w.Write(data)
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Error sending response"))
		return false
	}
	return true
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	error   error
}

func InternalServerError(l *zerolog.Logger, w http.ResponseWriter, e error) {
	er := ErrorResponse{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	WriteError(l, w, er)
}

func NotFoundError(l *zerolog.Logger, w http.ResponseWriter, msg string, e error) {
	er := ErrorResponse{Status: http.StatusNotFound, Message: msg, error: e}
	WriteError(l, w, er)
}

func BadRequestError(l *zerolog.Logger, w http.ResponseWriter, msg string, e error) {
	er := ErrorResponse{Status: http.StatusBadRequest, Message: msg, error: e}
	WriteError(l, w, er)
}

func WriteError(l *zerolog.Logger, w http.ResponseWriter, er ErrorResponse) {
	if er.error != nil {
		if er.Status > 500 {
			l.Error().Msg(er.error.Error())
		} else {
			l.Info().Msg(er.error.Error())
		}
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(er.Status)
	err := json.NewEncoder(w).Encode(er)
	if err != nil {
		// cant do anything about it
		l.Error().Msgf("Error sending e rror response: %s\n", err)
	}
}
