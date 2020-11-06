package handlers

import (
	"net/http"

	"github.com/mat-greaves/test-api/internal/middleware"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	er := ErrorResponse{
		Status:  http.StatusNotFound,
		Message: "Not Found",
	}
	WriteJSON(l, w, http.StatusNotFound, er)
}
