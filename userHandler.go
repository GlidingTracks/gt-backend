package main

import (
	"context"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/GlidingTracks/gt-backend/constant"
	model "github.com/GlidingTracks/gt-backend/models"
)

// CreateNewUser creates an user into the firebase user system
// Either returns the created users uId or an error.
func CreateNewUser(app *firebase.App, u model.User) (string, error) {
	if app == nil {
		return "", errors.New(constant.ErrorCouldNotContactFirebase)
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return "", err
	}

	params := (&auth.UserToCreate{}).
		Email(u.Email).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	user, err := client.CreateUser(ctx, params)

	if err != nil {
		return "", err
	}

	return user.UID, nil
}

// UpdateUser updates a user, must provide an uid in the u param. All fields are used, so if some is left out they will return to
// nil or default value.
func UpdateUser(app *firebase.App, u model.User) (model.User, error) {
	var uu model.User

	if app == nil {
		return uu, errors.New(constant.ErrorCouldNotContactFirebase)
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return uu, err
	}

	params := (&auth.UserToUpdate{}).
		Email(u.Email).
		EmailVerified(u.EmailVerified).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	ur, err := client.UpdateUser(ctx, u.Uid, params)

	if err != nil {
		return uu, err
	}

	uu = convertFromUserRecordToUser(*ur)

	return uu, nil
}

// GetUser fetches a user from firebase based on uId.
func GetUser(app *firebase.App, uID string) (model.User, error) {
	var u model.User

	if app == nil {
		return u, errors.New(constant.ErrorCouldNotContactFirebase)
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return u, err
	}

	ur, err := client.GetUser(ctx, uID)

	if err != nil {
		return u, err
	}

	u = convertFromUserRecordToUser(*ur)

	return u, nil
}

// DeleteUser will delete a user based on it's id from firebase.
func DeleteUser(app *firebase.App, uID string) error {
	if app == nil {
		return errors.New(constant.ErrorCouldNotContactFirebase)
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}

	err = client.DeleteUser(ctx, uID)
	if err != nil {
		return err
	}

	return nil
}

func checkIfUserExcists() error {
	return errors.New(constant.ErrorNotImplemented)
}

// convertFromUserRecordToUser takes a UserRecord, which a class from firebase, but it contains a lot of other uninteresting
// metadata, so as of now we stick to our model and just convert.
func convertFromUserRecordToUser(ur auth.UserRecord) model.User {
	var u model.User

	u.PhoneNumber = ur.PhoneNumber
	u.DisplayName = ur.DisplayName
	u.Email = ur.Email
	u.EmailVerified = ur.EmailVerified
	u.Uid = ur.UID
	u.Disabled = ur.Disabled

	return u
}
