package gtbackend

import (
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityMiddleware_CheckIncomingRequests(t *testing.T) {
	app, token := RetrieveFirebaseIDToken()

	sec := SecurityMiddleware{App: app}

	server := mux.NewRouter()
	server.Use(sec.CheckIncomingRequests)

	req, err := http.NewRequest("GET", "/getTracks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", token)

	rr := httptest.NewRecorder()
	server.HandleFunc("/getTracks", MockHandler)

	server.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Error("Wrong code returned")
	}
}

type authResponse struct {
	Kind         string
	IdToken      string
	RefreshToken string
	ExpiresIn    string
}

func RetrieveFirebaseIDToken() (app *firebase.App, token string) {
	config := &firebase.Config{
		StorageBucket: "gt-backend-8b9c2.appspot.com",
	}
	opt := option.WithCredentialsFile(constant.GoogleServiceCredName)

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

	token = resParsed.IdToken

	return
}
