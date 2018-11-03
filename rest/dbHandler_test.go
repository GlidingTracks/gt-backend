package rest

import (
	"encoding/json"
	"github.com/GlidingTracks/gt-backend/models"
	"github.com/GlidingTracks/gt-backend/testutils"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
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
	app, token := testutils.RetrieveFirebaseIDToken()
	values := map[string]io.Reader{
		"file": mustOpen("../testdata/testIgc.igc"),
	}

	r := CompleteRouterSetup(app)

	// Set up insertTrack
	req, err := createMultipart(values, "/insertTrack", "POST")
	if err != nil {
		t.Error("Could not create multipart")
	}
	req.Header.Set("token", token)

	// Run insertTrack
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Result().StatusCode != 200 {
		t.Error("Failed InsertTrack")
	}

	// Extract data from insertTrack
	insertRwBody, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Error("Failed reading body of InsertTrack")
	}
	var insertBody models.IgcMetadata
	err = json.Unmarshal(insertRwBody, &insertBody)
	if err != nil {
		t.Error("Failed extracting metadata response of InsertTrack")
	}
	// insertTrack DONE

	// Set up getTracks
	req = httptest.NewRequest("GET", "/getTracks", nil)
	req.Header.Set("token", token)
	req.Header.Set("queryType", "Private")
	req.Header.Set("orderDirection", "Desc")

	// Run getTracks
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Result().StatusCode != 200 {
		t.Error("Failed InsertTrack")
	}

	// Extract data from getTracks
	getTracksRwBody, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Error("Failed reading body of InsertTrack")
	}
	var getTracksBody []models.IgcMetadata
	err = json.Unmarshal(getTracksRwBody, &getTracksBody)
	if err != nil {
		t.Error("Failed extracting metadata response of InsertTrack")
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
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Result().StatusCode != 200 {
		t.Error("Failed GetTrack")
	}

	// Extract data from getTrack
	getTrackBody, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Error("Failed GetTrack body read")
	}
	var builder strings.Builder
	builder.Write(getTrackBody)
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

	// Run deleteTrack
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Result().StatusCode != 200 {
		t.Error("Failed GetTrack")
	}

	// Extract data from deleteTrack
	deleteTrackBody, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Error("Failed GetTrack body read")
	}
	builder.Write(deleteTrackBody)
	parsedDeleteTrackBody := builder.String()
	builder.Reset()

	// Test to see if response is equal to TrackID (then it was successful)
	if parsedDeleteTrackBody != insertBody.TrackID {
		t.Error("Failed DeleteTrack body read")
	}
}