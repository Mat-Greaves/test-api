package models

import (
	"context"
	"log"
	"os"
	"regexp"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

var us UserStorer
var collection *mongo.Collection

func TestMain(m *testing.M) {
	db, err := NewDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	collection = db.Database("test-api").Collection("users")
	us = NewUserStore(collection)
	os.Exit(m.Run())
}

func cleanup(t *testing.T) {
	if err := collection.Drop(context.Background()); err != nil {
		t.Fatalf("failed to cleanup database: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	t.Run("fail no age", func(t *testing.T) {
		cleanup(t)
		newUser := NewUser{
			Name: "Mat",
		}
		_, err := us.CreateUser(context.TODO(), &newUser)
		if err == nil {
			t.Errorf("Expected user create to fail")
		}
	})

	t.Run("fail no name", func(t *testing.T) {
		newUser := NewUser{
			Age: 31,
		}
		_, err := us.CreateUser(context.TODO(), &newUser)
		if err == nil {
			t.Errorf("Expected user create to fail")
		}
	})

	t.Run("create user", func(t *testing.T) {
		newUser := NewUser{
			Name: "Mat",
			Age:  31,
		}
		user, err := us.CreateUser(context.TODO(), &newUser)
		if err != nil {
			t.Errorf("Failed to create user: %s", err)
		}
		if user.Age != newUser.Age {
			t.Errorf("Age got: %d, want: %d", user.Age, newUser.Age)
		}
		if user.Name != newUser.Name {
			t.Errorf("Name got: %s, want: %s", user.Name, newUser.Name)
		}
		r := "[0-9a-fA-F]{24}"
		if matched, _ := regexp.Match(r, []byte(user.ID)); !matched {
			t.Errorf("Expected user ID to match %s", r)
		}
	})
}

func TestAllUsers(t *testing.T) {
	t.Run("no users", func(t *testing.T) {
		cleanup(t)
		users, err := us.AllUsers(context.Background())
		if err != nil {
			t.Errorf("AllUsers failed, expected success: %s", err)
		}
		if len(users) != 0 {
			t.Errorf("Got length %d, wanted %d", len(users), 0)
		}
	})

	t.Run("one user", func(t *testing.T) {
		cleanup(t)
		newUser := NewUser{
			Name: "Mat",
			Age:  31,
		}
		user, err := us.CreateUser(context.TODO(), &newUser)
		if err != nil {
			t.Fatalf("Setup failed, could not insert user: %s", err)
		}
		users, err := us.AllUsers(context.Background())
		if err != nil {
			t.Errorf("AllUsers failed, expected success: %s", err)
		}
		if len(users) != 1 {
			t.Errorf("Got length %d, wanted %d", len(users), 0)
		}
		if users[0] != user {
			t.Errorf("Did not get expected user, got %+v want %+v", users[0], user)
		}
	})
}
