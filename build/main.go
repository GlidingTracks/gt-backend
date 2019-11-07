package main

import (
	"errors"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/rest"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

const filename = "main.go"

// main is the first entry-point in application.
func main() {
	log := gtbackend.DebugLogPrepareHeader(filename, "main")
	app, err := initializeFirebase()
	if err != nil {
		gtbackend.LogFatalErrNoMsg(log, err)
	}

	rCors := rest.CompleteRouterSetup(app)

	port := os.Getenv("PORT")
	port = "8080"
	if port == "" {
		gtbackend.LogFatalNoErrMsg(log, "$PORT must be set")
	}

	err = http.ListenAndServe(":"+port, rCors)
	gtbackend.LogFatalErrNoMsg(log, err)
}

// initializeFirebase gets a App object from Firebase, based on the service account credentials.
func initializeFirebase() (app *firebase.App, err error) {
	log := gtbackend.DebugLogPrepareHeader(filename, "initializeFirebase")
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
		gtbackend.LogFatalErrMsg(log, err, "error initializing app")
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
