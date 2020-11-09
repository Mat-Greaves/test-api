package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mat-greaves/test-api/api"
	"github.com/mat-greaves/test-api/internal/rest/middleware"
	"github.com/mat-greaves/test-api/internal/service/users"
	"github.com/pkg/errors"
)

type createUserRequest struct {
	Name string `json:"name" validate:"required,lt=10"`
	Age  int    `json:"age" validate:"required,gte=0,lt=99"`
}

type UserHandler struct {
	validate *validator.Validate
	service  users.UserServicer
}

func NewUserHandler(validate *validator.Validate, service users.UserServicer) *UserHandler {
	return &UserHandler{
		validate: validate,
		service:  service,
	}
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	input := createUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		BadRequestError(l, w, err.Error(), errors.Wrap(err, "Error decoding json request"))
		return
	}
	err = uh.validate.Struct(input)
	if err != nil {
		BadRequestError(l, w, err.Error(), errors.Wrap(err, "Error validating user struct"))
		return
	}
	u := users.CreateUserRequest{
		Name: input.Name,
		Age:  input.Age,
	}
	user, err := uh.service.CreateUser(r.Context(), u)
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Failed to create user"))
		return
	}
	WriteJSON(l, w, http.StatusCreated, user)
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	u, err := uh.service.GetUsers(r.Context())
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Error retrieving users"))
		return
	}
	WriteJSON(l, w, http.StatusOK, convertUsers(u))
}

// convert internal User type to api User type for json tags
func convertUsers(u []users.User) []api.User {
	apiUsers := make([]api.User, len(u))
	for i, v := range u {
		apiUsers[i] = api.User(v)
	}
	return apiUsers
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	vars := mux.Vars(r)
	id := vars["id"]
	if err := uh.service.DeleteUser(r.Context(), id); err != nil {
		InternalServerError(l, w, errors.Wrapf(err, "Error deleting user id: %s", id))
		return
	}
	WriteJSON(l, w, http.StatusNoContent, nil)
}
