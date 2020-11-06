package models

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// these are the methods our user store supports. The userStore implementation
// is on particular implementation. This allows us to replace this implementation
// during our tests with a simple interface rather than using a concrete type that relies
// on the database reference
type UserStorer interface {
	CreateUser(context.Context, *NewUser) (User, error)
	AllUsers(context.Context) ([]User, error)
	DeleteUser(context.Context, string) error
}

// don't export type, everything external only needs to know about the
// UserStorer interface
type userStore struct {
	db *mongo.Collection
}

func NewUserStore(db *mongo.Collection) UserStorer {
	return &userStore{
		db: db,
	}
}

type NewUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	ID   string `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" `
	Age  int    `json:"age" `
}

func (us *userStore) CreateUser(ctx context.Context, user *NewUser) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("expected user name to be set")
	}
	if user.Age == 0 {
		return User{}, errors.New("expected user to have an age")
	}
	res, err := us.db.InsertOne(ctx, user)
	if err != nil {
		return User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user := User{
			ID:   oid.Hex(),
			Name: user.Name,
			Age:  user.Age,
		}
		return user, nil
	} else {
		return User{}, errors.New("Inserted id not json, bad?")
	}
}

func (us *userStore) AllUsers(ctx context.Context) ([]User, error) {
	var users []User
	cursor, err := us.db.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find users")
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &users); err != nil {

		return nil, errors.Wrap(err, "failed to read cursor")
	}
	return users, nil
}

func (us *userStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result := us.db.FindOneAndDelete(ctx, bson.M{"_id": oid})
	return result.Err()
}
