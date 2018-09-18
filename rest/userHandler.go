// Package rest contains everything related to the public RESTapi.
// Everything from handling the routes and performing logic
package rest

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	model "github.com/GlidingTracks/gt-backend/models"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

// Used in debugging. TODO remove before prod
const fileName = "userHandler.go"

// UserHandler handler object for acting upon a User in firebase.
// Contains a Context var as well as the different routes.
type UserHandler struct {
	Ctx            Context
	CreateUserPage string
	UpdateUserPage string
	DeleteUserPage string
	GetUserPage    string
}

// Bind sets up the routes to the mux router.
func (userHandler UserHandler) Bind(r *mux.Router) {
	r.HandleFunc("/createUser", userHandler.createUserPage).Methods(constant.Post)
	r.HandleFunc("/updateUser", userHandler.updateUserPage).Methods(constant.Post)
	r.HandleFunc("/deleteUser", userHandler.deleteUserPage).Queries("uId", "{uId}")
	r.HandleFunc("/getUser", userHandler.getUserPage).Queries("uId", "{uId}")
}

// createUserPage is the handler for creating users.
func (userHandler UserHandler) createUserPage(w http.ResponseWriter, r *http.Request) {
	var u model.User

	err := unmarshal(r, &u)
	if err != nil {
		logrus.Error(err)
		http.Error(w, errors.New(constant.ErrorProcessBodyFailed).Error(), http.StatusBadRequest)
		return
	}

	// Try to create a user in firebase
	_, err = createNewUser(userHandler.Ctx.App, u)
	if err != nil {
		logrus.Error(err)
		http.Error(w, errors.New(constant.ErrorCouldNotCreateUser).Error(), http.StatusBadRequest)
	}
}

// updateUserPage is the handler for updating a user.
func (userHandler UserHandler) updateUserPage(w http.ResponseWriter, r *http.Request) {
	var u model.User

	err := unmarshal(r, &u)
	if err != nil {
		logrus.Error(err)
		http.Error(w, errors.New(constant.ErrorProcessBodyFailed).Error(), http.StatusBadRequest)
		return
	}

	_, err = updateUser(userHandler.Ctx.App, u)
	if err != nil {
		logrus.Error(err)
		http.Error(w, errors.New(constant.ErrorCouldNotUpdateUser).Error(), http.StatusBadRequest)
	}
}

// deleteUserPage is the handler for deleting a user.
func (userHandler UserHandler) deleteUserPage(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	if queries == nil {
		http.Error(w, errors.New(constant.ErrorNoUIDProvided).Error(), http.StatusBadRequest)
		return
	}

	uID := queries.Get("uId")
	if uID == "" {
		http.Error(w, errors.New(constant.ErrorNoUIDProvided).Error(), http.StatusBadRequest)
		return
	}

	err := deleteUser(userHandler.Ctx.App, uID)
	if err != nil {
		http.Error(w, errors.New(constant.ErrorDeleteUser).Error(), http.StatusBadRequest)
		return
	}
}

// getUserPage is the handler for fetching a user from firebase.
func (userHandler UserHandler) getUserPage(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	if queries == nil {
		logrus.Error("no queries detected in path")
		http.Error(w, errors.New(constant.ErrorNoUIDProvided).Error(), http.StatusBadRequest)
		return
	}

	uID := queries.Get("uId")
	if uID == "" {
		logrus.Error("no uId specified in getUserPage")
		http.Error(w, errors.New(constant.ErrorNoUIDProvided).Error(), http.StatusBadRequest)
		return
	}

	u, err := getUser(userHandler.Ctx.App, uID)
	if err != nil {
		gtbackend.DebugLog(fileName, "getUserPage", err)
		http.Error(w, errors.New(constant.ErrorCouldNotGetUser).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

// createNewUser creates an user into the firebase user system.
// Either returns the created users uId or an error.
func createNewUser(app *firebase.App, u model.User) (uID string, err error) {
	if app == nil {
		err = errors.New(constant.ErrorCouldNotContactFirebase)
		return
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return
	}

	params := (&auth.UserToCreate{}).
		Email(u.Email).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	user, err := client.CreateUser(ctx, params)

	if err != nil {
		return
	}

	uID = user.UID
	return
}

// updateUser updates a user, must provide an uid in the u param. All fields are used, so if some is left out they will return to
// nil or default value.
func updateUser(app *firebase.App, u model.User) (user model.User, err error) {
	if app == nil {
		err = errors.New(constant.ErrorCouldNotContactFirebase)
		return
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return
	}

	params := (&auth.UserToUpdate{}).
		Email(u.Email).
		EmailVerified(u.EmailVerified).
		PhoneNumber(u.PhoneNumber).
		Password(u.Password).
		DisplayName(u.DisplayName).
		Disabled(u.Disabled)

	ur, err := client.UpdateUser(ctx, u.UID, params)
	if err != nil {
		return
	}

	user = convertFromUserRecordToUser(*ur)

	return
}

// getUser fetches a user from firebase based on uId.
func getUser(app *firebase.App, uID string) (user model.User, err error) {
	if app == nil {
		err = errors.New(constant.ErrorCouldNotContactFirebase)
		return
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return
	}

	ur, err := client.GetUser(ctx, uID)
	if err != nil {
		return
	}

	user = convertFromUserRecordToUser(*ur)

	return
}

// deleteUser will delete a user based on it's id from firebase.
func deleteUser(app *firebase.App, uID string) (err error) {
	if app == nil {
		err = errors.New(constant.ErrorCouldNotContactFirebase)
		return
	}

	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		return
	}

	err = client.DeleteUser(ctx, uID)
	if err != nil {
		return
	}

	return
}

func checkIfUserExists() error {
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
	u.UID = ur.UID
	u.Disabled = ur.Disabled
	u.PhoneNumber = ur.PhoneNumber

	return u
}

// unmarshal takes a request and and interface. Tries to decode r.Body into the interface.
func unmarshal(r *http.Request, v interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&v); err != nil {
		return
	}

	defer r.Body.Close()

	return
}
