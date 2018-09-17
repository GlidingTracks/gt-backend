package main

import (
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend/models"
	"google.golang.org/api/option"
	"log"
	"testing"
)

func initApp() (app *firebase.App) {
	opt := option.WithCredentialsFile("gt-backend-8b9c2-firebase-adminsdk-0t965-d5b53ac637.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}

func TestCreateNewUser(t *testing.T) {
	var userID string

	t.Run("Create user", func(t *testing.T) {
		var u = new(models.User)
		u.Email = "test@test.com"
		u.PhoneNumber = "+4799999999"
		u.DisplayName = "test"
		u.Password = "test1234"

		app := initApp()

		uID, err := CreateNewUser(app, *u)

		if err != nil {
			t.Error("Could not create user", err)
		}

		if uID == "" {
			t.Error("Could not create user")
		}

		userID = uID
	})

	t.Run("Get", func(t *testing.T) {
		app := initApp()

		_, err := GetUser(app, userID)
		if err != nil {
			t.Error("Could not get user", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		var u = new(models.User)
		u.Uid = userID
		u.Email = "test@test.com"
		u.PhoneNumber = "+4799999999"
		u.DisplayName = "testUpdate"
		u.Password = "test1234"

		app := initApp()

		uu, err := UpdateUser(app, *u)
		if err != nil {
			t.Error("Could not update user", err)
		}

		if uu.DisplayName != "testUpdate" {
			t.Error("User not updated correctly")
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		if userID == "" {
			t.Skip("Previous test failed")
		}

		app := initApp()

		err := DeleteUser(app, userID)
		if err != nil {
			t.Error("Could not delete user", err)
		}
	})
}
