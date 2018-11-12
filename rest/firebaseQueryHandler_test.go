package rest

import (
	"github.com/GlidingTracks/gt-backend/constant"
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

func TestUpdatePrivacy(t *testing.T) {
	app := testutils.InitializeFirebaseTest()

	trackID := "scf6Xw4pwCKGeLjrVJHo"

	// Test setting privacy to TRUE
	data, err := UpdatePrivacy(app, trackID, constant.TestUID, true)
	if err != nil {
		t.Error("UpdatePrivacy failed updating setting")
	}
	if data.Privacy != true {
		t.Error("Privacy should be TRUE")
	}

	// Test setting privacy to FALSE
	data, err = UpdatePrivacy(app, trackID, constant.TestUID, false)
	if err != nil {
		t.Error("UpdatePrivacy failed updating setting")
	}
	if data.Privacy != false {
		t.Error("Privacy should be FALSE")
	}

	// Test setting with wrong UID, should fail
	data, err = UpdatePrivacy(app, trackID, "TotallyLegitUID", false)
	if err == nil {
		t.Error("This should actually have an error")
	}
}

func TestInsertTrackPoint(t *testing.T) {
	app := testutils.InitializeFirebaseTest()

	trackID := "scf6Xw4pwCKGeLjrVJHo"
	var testTrackPointArray []models.TrackPoint
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData)

	data, err := InsertTrackPoint(app, trackID, constant.TestUID, testTrackPointArray)
	if err != nil {
		t.Error("InsertTrackPoint failed uploading track points")
	}
	if data.TrackPoints[0] != InsertTrackPointTestData {
		t.Error("InsertTrackPoint insertion should be same as object that was sent")
	}
}

// DeleteTrack method tested in fileUploadHandler_test
// TakeOwnership method tested in fileUploadHandler_test
