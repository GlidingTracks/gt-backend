package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/Sirupsen/logrus"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
)

// InitializeFirebaseTest Normal firebase initialization for testing purposes, no auth token
func InitializeFirebaseTest() (app *firebase.App) {
	app = InitializeFirebaseTestCredFile(true)
	return
}

// InitializeFirebaseTestCredFile Firebase initialization with auth token to get past security checking, with flag to open Firebase file in different folder
func InitializeFirebaseTestCredFile(credNotInFolder bool) (app *firebase.App) {
	config := &firebase.Config{
		StorageBucket: "gt-backend-8b9c2.appspot.com",
	}
	credPath := ""
	if credNotInFolder {
		credPath = "../" + constant.GoogleServiceCredName
	} else {
		credPath = constant.GoogleServiceCredName
	}
	opt := option.WithCredentialsFile(credPath)

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	return
}

type authResponse struct {
	Kind         string
	IDToken      string
	RefreshToken string
	ExpiresIn    string
}

// RetrieveFirebaseIDToken Firebase initialization with auth token to get past security checking
func RetrieveFirebaseIDToken() (app *firebase.App, token string) {
	app, token = RetrieveFirebaseIDTokenCredFile(true)
	return
}

// RetrieveFirebaseIDTokenCredFile Firebase initialization with auth token to get past security checking, with flag to open Firebase file in different folder
func RetrieveFirebaseIDTokenCredFile(credNotInFolder bool) (app *firebase.App, token string) {
	config := &firebase.Config{
		StorageBucket: "gt-backend-8b9c2.appspot.com",
	}

	credPath := ""
	if credNotInFolder {
		credPath = "../" + constant.GoogleServiceCredName
	} else {
		credPath = constant.GoogleServiceCredName
	}
	opt := option.WithCredentialsFile(credPath)

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		logrus.Fatalf("error initializing auth: %v\n", err)
	}

	token, err = client.CustomToken(context.Background(), "o1Sz791YSHby0PCe51JlxSD6G533")
	if err != nil {
		logrus.Fatalf("error setting custom token: %v\n", err)
	}

	var jsonStr = []byte(`{
	"token": "` + token + `",
	"returnSecureToken": true
	}`)
	res, err := http.Post("https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=AIzaSyAppC_L-VHnTM1ezOvuiVCoKfFzFu6f5ZU", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Printf("%+v\n", string(jsonStr))
		logrus.Fatalf("error retrieving id token: %v\n", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Fatalf("error reading body of id token retrieve req: %v\n", err)
	}

	var resParsed authResponse
	err = json.Unmarshal(resBody, &resParsed)
	if err != nil {
		logrus.Fatalf("error parsing json of id token request: %v\n", err)
	}

	token = resParsed.IDToken

	return
}
