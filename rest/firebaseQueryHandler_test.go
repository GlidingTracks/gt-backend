package rest

import (
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/GlidingTracks/gt-backend/testutils"
	"testing"
)

func TestGetTracks(t *testing.T) {
	app := testutils.InitializeFirebaseTest()

	testUID := "iP1dgAHJ2JNce4hGr9H0RugkCHP2"
	privateQ := models.NewFirebaseQuery(testUID, "1", "Private", "Asc")

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

	publicQ := models.NewFirebaseQuery(testUID, "1", "Public", "Asc")

	res, err = GetTracks(app, publicQ)
	if err != nil {
		t.Error("GetTracks failed public test!")
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

func TestGetTrack(t *testing.T) {
	app := testutils.InitializeFirebaseTest()

	data, err := GetTrack(app, "HAGOdywD9rQayoOOIHyd")
	if err != nil && len(data) < 1 {
		t.Error("Did not receive any data, should receive data", err)
	}
}
