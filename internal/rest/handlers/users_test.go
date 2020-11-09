package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/mat-greaves/test-api/internal/rest/middleware"
	"github.com/mat-greaves/test-api/internal/service/users"
	"github.com/rs/zerolog"
)

var uh *UserHandler

// test user store interface that returns fake data
type testUserService struct{}

// logger to use to inject into all of our handlers
var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func (*testUserService) CreateUser(ctx context.Context, user users.CreateUserRequest) (users.User, error) {
	return users.User{
		ID:   "testId",
		Name: user.Name,
		Age:  user.Age,
	}, nil
}

func (*testUserService) GetUsers(ctx context.Context) ([]users.User, error) {
	return []users.User{
		{
			ID:   "testId",
			Name: "test",
			Age:  50,
		},
	}, nil
}

func (*testUserService) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func TestMain(m *testing.M) {
	validate := validator.New()
	uh = NewUserHandler(validate, &testUserService{})
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	// create a context with logger usually populated by middleware
	ctx := middleware.GetContextWithLogger(context.Background(), &logger, "test123")
	type createUserTest struct {
		name       string
		method     string
		input      string
		want       string
		statusCode int
	}
	tt := []createUserTest{
		{
			name:       "no name",
			method:     "POST",
			input:      `{ "age": 30 }`,
			want:       `{"status":400,"message":"Key: 'createUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			statusCode: 400,
		},
		{
			name:       "no age",
			method:     "POST",
			input:      `{ "name": "Mat" }`,
			want:       `{"status":400,"message":"Key: 'createUserRequest.Age' Error:Field validation for 'Age' failed on the 'required' tag"}`,
			statusCode: 400,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest(tc.method, "/users", bytes.NewBuffer([]byte(tc.input)))
			if err != nil {
				t.Fatal(err)
			}
			r.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()
			uh.CreateUser(rr, r.WithContext(ctx))
			if status := rr.Code; status != tc.statusCode {
				t.Errorf("got status: %d, wanted %d", rr.Code, tc.statusCode)
			}
			if body := strings.TrimSpace(rr.Body.String()); body != tc.want {
				t.Errorf("want '%s' got '%s'", tc.want, body)
			}
		})
	}
}
