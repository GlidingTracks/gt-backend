package rest

import (
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"net/http/httptest"
	"testing"
)

func TestDbHandler_Implementations(t *testing.T) {
	var handler interface{} = &DbHandler{}
	if _, implemented := handler.(MuxRouteBinder); !implemented {
		t.Error("does not implement MuxRouteBinder")
	}
}

func InitializeFirebaseTest() (app *firebase.App) {
	config := &firebase.Config{
		StorageBucket: "gt-backend-8b9c2.appspot.com",
	}
	opt := option.WithCredentialsFile("../" + constant.GoogleServiceCredName)

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	return
}

func TestDbHandler(t *testing.T) {
	testRouter := mux.NewRouter()

	dbHandler := DbHandler{
		Context{},
		"",
		"",
		"",
		"",
	}

	t.Run("Insert", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req := httptest.NewRequest("GET", "/getTracks", nil)
		testRouter.HandleFunc("/getTracks", dbHandler.getTracksPage)
		testRouter.ServeHTTP(rr, req)

		if rr.Code != 400 {
			t.Error("Expected error")
		}
	})

}
