package main

import (
	"encoding/json"
	"errors"
	"firebase.google.com/go"
	"fmt"
	"github.com/GlidingTracks/gt-backend/constant"
	model "github.com/GlidingTracks/gt-backend/models"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

// main is the first entry-point in application.
func main() {
	initializeFirebase()

	r := mux.NewRouter()

	r.HandleFunc("/", startPage)
	r.HandleFunc("/createUser", createUserPage).Methods("POST")
	r.HandleFunc("/UpdateUser", updateUserPage).Methods("POST")
	r.HandleFunc("/deleteUser", deleteUserPage).Queries("uId", "{uId}")
	r.HandleFunc("/getUser", getUserPage).Queries("uId", "{uId}")

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

// startPage redirects every non-existing path to url: localhost:8080/.
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go - Gliding tracks\n")
}

// createUserPage is the endpoint for creating users
func createUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	var u model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	_, err := CreateNewUser(app, u)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(400)
	}
}

// updateUserPage is the endpoint for updating a user
func updateUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	var u model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	_, err := UpdateUser(app, u)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(400)
	}
}

// deleteUserPage is the endpoint for deleting a user.
func deleteUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	queries := r.URL.Query()
	if queries == nil {
		http.Error(w, errors.New(constant.ErrorNoUidProvided).Error(), http.StatusBadRequest)
		return
	}

	uID := queries.Get("uId")
	if uID == "" {
		http.Error(w, errors.New(constant.ErrorNoUidProvided).Error(), http.StatusBadRequest)
		return
	}

	err := DeleteUser(app, uID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// getUserPage is the endpoint for fetching a user from firebase.
func getUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	queries := r.URL.Query()
	if queries == nil {
		http.Error(w, errors.New(constant.ErrorNoUidProvided).Error(), http.StatusBadRequest)
		return
	}

	uID := queries.Get("uId")
	if uID == "" {
		http.Error(w, errors.New(constant.ErrorNoUidProvided).Error(), http.StatusBadRequest)
		return
	}

	u, err := GetUser(app, uID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

// initializeFirebase gets a App object from Firebase, based on the service account credentials.
func initializeFirebase() *firebase.App {
	opt := option.WithCredentialsFile(constant.GoogleServiceCredName)
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
