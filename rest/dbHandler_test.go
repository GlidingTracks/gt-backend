package rest

import (
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/Sirupsen/logrus"
	"google.golang.org/api/option"
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
