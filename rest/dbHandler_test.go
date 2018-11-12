package rest

import (
	"encoding/json"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/GlidingTracks/gt-backend/testutils"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDbHandler_Implementations(t *testing.T) {
	var handler interface{} = &DbHandler{}
	if _, implemented := handler.(MuxRouteBinder); !implemented {
		t.Error("does not implement MuxRouteBinder")
	}
}

func TestDbHandler(t *testing.T) {
	testRouter := mux.NewRouter()

	dbHandler := DbHandler{
		Context{},
		"",
		"",
		"",
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

// E2E test of all functions in DbHandler executed in following order:
// InsertTrack -> GetTracks -> GetTrack -> DeleteTrack
func TestIntegratedDbHandlerTest(t *testing.T) {
	app := testutils.InitializeFirebaseTest()
	scraperToken := testutils.RetrieveFirebaseIDToken(app, constant.ScraperUID)
	token := testutils.RetrieveFirebaseIDToken(app, constant.TestUID)
	values := map[string]io.Reader{
		"file":    mustOpen("../testdata/testIgc.igc"),
		"private": strings.NewReader("false"),
	}

	r := CompleteRouterSetup(app)

	// Set up insertTrack
	req, err := createMultipart(values, "/insertTrack", "POST")
	if err != nil {
		t.Error("Could not create multipart")
	}
	req.Header.Set("token", scraperToken)

	// Run insertTrack
	ret := testutils.TestRoute(req, r, "InsertTrack", t, http.StatusOK)
	var insertBody models.IgcMetadata
	err = json.Unmarshal(ret, &insertBody)
	if err != nil {
		t.Error("Failed extracting metadata response of InsertTrack")
	}
	// insertTrack DONE

	// Set up takeOwnership
	req = httptest.NewRequest("PUT", "/takeOwnership", nil)
	req.Header.Set("token", token)
	req.Header.Set("trackID", insertBody.TrackID)
	if insertBody.UID != constant.ScraperUID {
		t.Error("Track should be owned by the scraper before taking ownership")
	}

	// Run TakeOwnership
	ret = testutils.TestRoute(req, r, "TakeOwnership", t, http.StatusOK)
	var takeOwnershipBody models.IgcMetadata
	err = json.Unmarshal(ret, &takeOwnershipBody)
	if err != nil {
		t.Error("Failed extracting metadata response of TakeOwnership")
	}

	if takeOwnershipBody.UID != constant.TestUID {
		t.Error("Track should now be owned by the TestUID")
	}

	// takeOwnership DONE

	// Set up updatePrivacy
	req = httptest.NewRequest("PUT", "/updatePrivacy", nil)
	req.Header.Set("token", token)
	req.Header.Set("trackID", insertBody.TrackID)
	req.Header.Set("private", "true")
	if insertBody.Privacy != false {
		t.Error("Privacy should be false before updating privacy to true")
	}

	// Run UpdatePrivacy
	ret = testutils.TestRoute(req, r, "UpdatePrivacy", t, http.StatusOK)
	var updatePrivacyBody models.IgcMetadata
	err = json.Unmarshal(ret, &updatePrivacyBody)
	if err != nil {
		t.Error("Failed extracting metadata response of UpdatePrivacy")
	}

	if updatePrivacyBody.Privacy != true {
		t.Error("Privacy setting should be changed to true")
	}
	// updatePrivacy DONE

	// Set up insertTrackPoint
	// Set up the object to send in correct format
	var trackPoints []models.TrackPoint
	trackPoints = append(trackPoints, testutils.InsertTrackPointTestData)
	trackPoints = append(trackPoints, testutils.InsertTrackPointTestData)
	trackPoints = append(trackPoints, testutils.InsertTrackPointTestData)
	trackPoints = append(trackPoints, testutils.InsertTrackPointTestData)
	trackPointsJson, err := json.Marshal(trackPoints)
	if err != nil {
		t.Error("Error parsing JSON of insertTrackPoint")
	}
	var tempBuilder strings.Builder
	tempBuilder.Write(trackPointsJson)
	trackPointsString := tempBuilder.String()

	// Set up the request
	req = httptest.NewRequest("PUT", "/insertTrackPoint", nil)
	req.Header.Set("token", token)
	req.Header.Set("trackID", insertBody.TrackID)
	req.Header.Set("trackPoints", trackPointsString)

	// insertTrackPoint DONE
	ret = testutils.TestRoute(req, r, "InsertTrackPoint", t, http.StatusOK)
	var insertTracksPointBody models.IgcMetadata
	err = json.Unmarshal(ret, &insertTracksPointBody)
	if err != nil {
		t.Error("Failed extracting metadata response of InsertTrackPoint")
	}

	if insertTracksPointBody.TrackPoints[0] != testutils.InsertTrackPointTestData {
		t.Error("InsertTrackPoint insertion should be same as object that was sent")
	}

	// Set up getTracks
	req = httptest.NewRequest("GET", "/getTracks", nil)
	req.Header.Set("token", token)
	req.Header.Set("queryType", "Private")
	req.Header.Set("orderDirection", "Desc")

	// Run getTracks
	ret = testutils.TestRoute(req, r, "GetTracks", t, http.StatusOK)
	var getTracksBody []models.IgcMetadata
	err = json.Unmarshal(ret, &getTracksBody)
	if err != nil {
		t.Error("Failed extracting metadata response of GetTracks")
	}

	// Test getTracks data against insertTrack data
	// There should be at least one track in body, and the first needs to match uploaded track
	if insertBody.TrackID != getTracksBody[0].TrackID {
		t.Error("Expected inserted trackID to match 1st getTracks trackID")
	}
	// getTracks DONE

	// Set up getTrack
	req = httptest.NewRequest("GET", "/getTrack", nil)
	req.Header.Set("token", token)
	req.Header.Set("trackID", insertBody.TrackID)

	// Run getTrack
	ret = testutils.TestRoute(req, r, "GetTrack", t, http.StatusOK)
	var builder strings.Builder
	builder.Write(ret)
	parsedGetTrackBody := builder.String()
	builder.Reset()

	// Extract file to test against getTrack
	fileBody, err := ioutil.ReadFile("../testdata/testIgc.igc")
	if err != nil {
		t.Error("Failed file read for GetTrack body test")
	}
	builder.Write(fileBody)
	parsedFileBody := builder.String()
	builder.Reset()

	// Test getTrack body against file that was uploaded
	if parsedGetTrackBody != parsedFileBody {
		t.Error("Failed GetTrack comparison")
	}
	// getTrack DONE

	// Set up deleteTrack
	req = httptest.NewRequest("DELETE", "/deleteTrack", nil)
	req.Header.Set("token", token)
	req.Header.Set("trackID", insertBody.TrackID)

	ret = testutils.TestRoute(req, r, "DeleteTrack", t, http.StatusOK)
	builder.Write(ret)
	parsedDeleteTrackBody := builder.String()
	builder.Reset()

	// Test to see if response is equal to TrackID (then it was successful)
	if parsedDeleteTrackBody != insertBody.TrackID {
		t.Error("Failed DeleteTrack body read")
	}
}
