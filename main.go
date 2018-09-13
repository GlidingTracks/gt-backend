package main

import (
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	model "github.com/GlidingTracks/gt-backend/models"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func main() {
	initializeFirebase()

	r := mux.NewRouter()

	r.HandleFunc("/", startPage)
	r.HandleFunc("/createUser", createUserPage).Methods("POST")
	r.HandleFunc("/updateUser", updateUserPage).Methods("POST")

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

// Redirect here is url: localhost:8080 is supplied
func startPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from go - Gliding tracks\n")
}

// Endpoint for creating users
func createUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	var u model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	err := createNewUser(app, u)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(400)
	}
}

// Endpoint for updating a user
func updateUserPage(w http.ResponseWriter, r *http.Request) {
	app := initializeFirebase()

	var u model.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		w.WriteHeader(400)
		return
	}

	defer r.Body.Close()

	err := updateUser(app, u)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(400)
	}
}

// Get a App object from Firebase, based on the service account credentials
func initializeFirebase() *firebase.App {
	opt := option.WithCredentialsFile("gt-backend-8b9c2-firebase-adminsdk-0t965-d5b53ac637.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
