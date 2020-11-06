package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mat-greaves/test-api/internal/middleware"
	"github.com/mat-greaves/test-api/internal/models"
	"github.com/pkg/errors"
)

type createUserRequest struct {
	Name string `json:"name" validate:"required,lt=10"`
	Age  int    `json:"age" validate:"required,gte=0,lt=99"`
}

type UserHandler struct {
	validate  *validator.Validate
	userStore models.UserStorer
}

func NewUserHandler(validate *validator.Validate, userStore models.UserStorer) *UserHandler {
	return &UserHandler{
		validate:  validate,
		userStore: userStore,
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
	u := models.NewUser{
		Name: input.Name,
		Age:  input.Age,
	}
	user, err := uh.userStore.CreateUser(r.Context(), &u)
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Failed to create user"))
		return
	}
	WriteJSON(l, w, http.StatusCreated, user)
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	users, err := uh.userStore.AllUsers(r.Context())
	if err != nil {
		InternalServerError(l, w, errors.Wrap(err, "Error retrieving users"))
		return
	}
	WriteJSON(l, w, http.StatusOK, users)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	l := middleware.GetLogger(r)
	vars := mux.Vars(r)
	id := vars["id"]
	if err := uh.userStore.DeleteUser(r.Context(), id); err != nil {
		// check if not found error
		// TODO: what's the better way to do this? https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
		matched, _ := regexp.Match(`no documents in result`, []byte(err.Error()))
		if matched {
			NotFoundError(l, w, "Not Found", err)
			return
		}
		// some kind of internal server error
		l.Error().Msgf("Error deleting user id: %s, error: %s", id, err)
		InternalServerError(l, w, err)
		return
	}
	WriteJSON(l, w, http.StatusNoContent, nil)
}
