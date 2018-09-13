package main

import (
	"context"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	model "github.com/GlidingTracks/gt-backend/models"
)

// Creates an user into the firebase user system
func createNewUser(app *firebase.App, u model.User) error {
	ctx := context.Background()

	client, err := app.Auth(ctx)

	params := (&auth.UserToCreate{}).
		Email(u.Email).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	_, err = client.CreateUser(ctx, params)

	if err != nil {
		return err
	}

	return nil
}

// Updates a user, must provide an uid in the u param. All fields are used, so if some is left out they will return to
// nil or default value
func updateUser(app* firebase.App, u model.User) error {
	ctx := context.Background()

	client, err := app.Auth(ctx)

	params := (&auth.UserToUpdate{}).
		Email(u.Email).
		EmailVerified(u.EmailVerified).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	_, err = client.UpdateUser(ctx, u.Uid, params)

	if err != nil {
		return err
	}

	return nil
}

func checkIfUserExcists() error {
	return errors.New("not implemented")
}
