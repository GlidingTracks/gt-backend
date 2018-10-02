package rest

import (
	"context"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/Sirupsen/logrus"
	"google.golang.org/api/option"
	"testing"
)

func TestGetTracks(t *testing.T) {
	opt := option.WithCredentialsFile("../" + constant.GoogleServiceCredName)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}

	testUID := "test2"
	privateQ := models.NewFirebaseQuery(testUID, "1", "Private", "Time", "Asc")

	res, err := GetTracks(app, privateQ)
	if err != nil {
		t.Error("GetTracks failed private test!")
	}
	if len(res) < 1 {
		t.Error("GetTracks private test with no results!")
	}
	for i := 0; len(res) < i; i++ {
		if res[i].UID != testUID {
			t.Error("GetTracks non-testUID data in private query!")
		}
	}

	publicQ := models.NewFirebaseQuery(testUID, "1", "Public", "Time", "Asc")

	res, err = GetTracks(app, publicQ)
	if err != nil {
		t.Error("GetTracks failed private test!")
	}
	if len(res) < 1 {
		t.Error("GetTracks public test with no results!")
	}
	for i := 0; len(res) < i; i++ {
		if res[i].UID == testUID {
			t.Error("GetTracks testUID data in public query!")
		}
	}
}
