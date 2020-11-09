package users

import (
	"context"

	"github.com/pkg/errors"
)

type UserService struct {
	store UserStorer
}

type UserServicer interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	DeleteUser(ctx context.Context, id string) error
}

// TODO: is service in the name here redundant?
func NewUserService(store UserStorer) UserServicer {
	return &UserService{
		store: store,
	}
}

type CreateUserRequest struct {
	Name string
	Age  int
}

// TODO: seems a bit gross having the store tags here too but saves all the conversion
type User struct {
	ID   string `bson:"_id,omitempty"`
	Name string
	Age  int
}

func (us *UserService) CreateUser(ctx context.Context, user CreateUserRequest) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("expected user name to be set")
	}
	if user.Age == 0 {
		return User{}, errors.New("expected user to have an age")
	}
	res, err := us.store.CreateUser(ctx, user)
	u := User{
		ID:   res.ID,
		Name: res.Name,
		Age:  res.Age,
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (us *UserService) GetUsers(ctx context.Context) ([]User, error) {
	users, err := us.store.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) DeleteUser(ctx context.Context, id string) error {
	if ok, err := us.store.DeleteUser(ctx, id); !ok {
		return errors.Wrap(err, "user service delete user error")
	}
	return nil
}
