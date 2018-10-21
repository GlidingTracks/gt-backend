package main

import (
	"errors"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/rest"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

// main is the first entry-point in application.
func main() {
	// TODO set correct level in prod
	logrus.SetLevel(logrus.DebugLevel)

	app, err := initializeFirebase()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	ctx := &rest.Context{
		App: app,
	}

	r := mux.NewRouter()

	userRoutes := &rest.UserHandler{
		Ctx:            *ctx,
		CreateUserPage: "/createUser",
		UpdateUserPage: "/updateUser",
		DeleteUserPage: "/deleteUser",
		GetUserPage:    "/getUser",
	}

	fileUploadRoutes := &rest.FileUploadHandler{
		Ctx:            *ctx,
		UploadFilePage: "/upload",
	}

	dbRoutes := &rest.DbHandler{
		Ctx:         *ctx,
		InsertTrack: "/insertTrack",
		GetTracks:   "/getTracks",
		GetTrack:    "/getTrack",
		DeleteTrack: "/deleteTrack",
	}

	userRoutes.Bind(r)
	fileUploadRoutes.Bind(r)
	dbRoutes.Bind(r)

	r.HandleFunc("/", startPage)

	port := os.Getenv("PORT")
	if port == "" {
		logrus.Fatal("$PORT must be set")
	}

	logrus.Fatal(http.ListenAndServe(":" + port, r))
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

	opt := option.WithCredentialsFile(constant.GoogleServiceCredName)

	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
		return
	}

	return
}

func checkIfFirebaseCredentialsExist() (exist bool) {
	exist = false
	_, err := os.Open(constant.GoogleServiceCredName)
	if err != nil {
		return
	}

	exist = true
	return
}

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

	_, err = f.WriteString(val)
	if err != nil {
		return
	}

	success = true
	return
}
