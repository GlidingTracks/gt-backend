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
)

// main is the first entry-point in application.
func main() {
	// TODO set correct level in prod
	logrus.SetLevel(logrus.DebugLevel)

	ctx := &rest.Context{
		App: initializeFirebase(),
	}

	r := mux.NewRouter()

	userRoutes := &rest.UserHandler{
		Ctx:            *ctx,
		CreateUserPage: "/createUser",
		UpdateUserPage: "/updateUser",
		DeleteUserPage: "/deleteUser",
		GetUserPage:    "/getUser",
	}

	userRoutes.Bind(r)

	fileUploadRoutes := &rest.FileUploadHandler{
		Ctx: *ctx,
		UploadFilePage: "/upload",
	}

	fileUploadRoutes.Bind(r)

	r.HandleFunc("/", startPage)

	logrus.Fatal(http.ListenAndServe(":8080", r))
}

// startPage redirects every non-existing path to url: localhost:8080/.
func startPage(w http.ResponseWriter, r *http.Request) {
	err := errors.New("page not found")
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// initializeFirebase gets a App object from Firebase, based on the service account credentials.
func initializeFirebase() (app *firebase.App) {
	opt := option.WithCredentialsFile(constant.GoogleServiceCredName)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	return
}
