package main

import (
	"errors"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/rest"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

// main is the first entry-point in application.
func main() {
	app, err := initializeFirebase()
	if err != nil {
		gtbackend.LogFatal(gtbackend.InternalLog{
			Err: err,
		})
	}

	ctx := &rest.Context{
		App: app,
	}

	r := mux.NewRouter()
	r.Use(gtbackend.LogIncomingRequests)

	userRoutes := &rest.UserHandler{
		Ctx:            *ctx,
		CreateUserPage: "/createUser",
		UpdateUserPage: "/updateUser",
		DeleteUserPage: "/deleteUser",
		GetUserPage:    "/getUser",
	}

	dbRoutes := &rest.DbHandler{
		Ctx:         *ctx,
		InsertTrack: "/insertTrack",
		GetTracks:   "/getTracks",
		GetTrack:    "/getTrack",
		DeleteTrack: "/deleteTrack",
	}

	userRoutes.Bind(r)
	dbRoutes.Bind(r)

	r.HandleFunc("/", startPage)

	port := os.Getenv("PORT")
	if port == "" {
		gtbackend.LogFatal(gtbackend.InternalLog{
			Msg: "$PORT must be set",
		})
	}

	gtbackend.LogFatal(gtbackend.InternalLog{
		Err: http.ListenAndServe(":"+port, r),
	})
}

// startPage redirects every non-existing path to url: localhost:8080/.
func startPage(w http.ResponseWriter, r *http.Request) {
	err := errors.New("page not found")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// initializeFirebase gets a App object from Firebase, based on the service account credentials.
func initializeFirebase() (app *firebase.App, err error) {
	if !checkIfFirebaseCredentialsExist() {
		if !tryCreateFirebaseCredsFromEnv() {
			err = errors.New("could not connect to DB")
			return
		}
	}

	config := &firebase.Config{
		StorageBucket: constant.FirebaseStorageBucket,
	}
	opt := option.WithCredentialsFile(constant.GoogleServiceCredName)

	app, err = firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		gtbackend.LogFatal(gtbackend.InternalLog{
			Msg: "error initializing app",
			Err: err,
		})
	}

	return
}

// checkIfFirebaseCredentialsExist will check for credential file, bool
func checkIfFirebaseCredentialsExist() (exist bool) {
	exist = false
	_, err := os.Open(constant.GoogleServiceCredName)
	if err != nil {
		return
	}

	exist = true
	return
}

// tryCreateFirebaseCredsFromEnv if cred content is loaded as a environment variable, create cred file from it
func tryCreateFirebaseCredsFromEnv() (success bool) {
	success = false

	val := os.Getenv(constant.GoogleCredEnvVar)
	if val == "" {
		return
	}

	f, err := os.Create(constant.GoogleServiceCredName)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.WriteString(val)
	if err != nil {
		return
	}

	success = true
	return
}
