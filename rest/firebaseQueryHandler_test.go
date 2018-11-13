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

	data, err := GetTrack(app, "HAGOdywD9rQayoOOIHyd", "123")
	if err != nil && len(data) < 1 {
		t.Error("Did not receive any data, should receive data", err)
	}

	_, err = GetTrack(app, "HAGOdywD9rQayoOOIHyd", "456")
	if err == nil {
		t.Error("Getting with wrong UID should throw error!", err)
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
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData1)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData2)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData3)
	testTrackPointArray = append(testTrackPointArray, InsertTrackPointTestData4)

	data, err := InsertTrackPoint(app, trackID, constant.TestUID, testTrackPointArray)
	if err != nil {
		t.Error("InsertTrackPoint failed uploading track points on TestUID")
	}
	if data.TrackPoints[0] != InsertTrackPointTestData1 {
		t.Error("InsertTrackPoint Object 0 should match object order of appending")
	}
	if data.TrackPoints[1] != InsertTrackPointTestData2 {
		t.Error("InsertTrackPoint Object 1 should match object order of appending")
	}
	if data.TrackPoints[2] != InsertTrackPointTestData3 {
		t.Error("InsertTrackPoint Object 2 should match object order of appending")
	}
	if data.TrackPoints[3] != InsertTrackPointTestData4 {
		t.Error("InsertTrackPoint Object 3 should match object order of appending")
	}

	data, err = InsertTrackPoint(app, trackID, "DummyID", testTrackPointArray)
	if err != nil {
		t.Error("InsertTrackPoint failed uploading track points on dummy ID")
	}

}

// DeleteTrack method tested in fileUploadHandler_test
// TakeOwnership method tested in fileUploadHandler_test
