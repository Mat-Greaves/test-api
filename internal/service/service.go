package service

import (
	"github.com/mat-greaves/test-api/internal/service/users"
)

type Service struct {
	Users users.UserServicer
}
