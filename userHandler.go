package main

import (
	"context"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	model "github.com/GlidingTracks/gt-backend/datastructures"
)

// Creates an user into the firebase user system
func createNewUser(app *firebase.App, u model.User) error {
	ctx := context.Background()

	client, err := app.Auth(ctx)

	params := (&auth.UserToCreate{}).
		Email(u.Email).
		Password("secretPassword").
		DisplayName(u.DisplayName).
		Disabled(false)

	_, err = client.CreateUser(ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func checkIfUserExcists() error {
	return errors.New("not implemented")
}
