package rest

import (
	"bytes"
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/testutils"
	"github.com/Sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProcessUploadRequestWrongContentType(t *testing.T) {
	app := testutils.InitializeFirebaseTest()
	values := map[string]io.Reader{
		"file": mustOpen("../testdata/text.txt"),
	}

	req, err := createMultipart(values, "/upload", "POST")
	req.Header.Set("uid", "123")
	if err != nil {
		t.Error("Could not create multipart")
	}

	code, _, err := ProcessUploadRequest(app, req)

	if err == nil && code != 415 {
		t.Error("Wrong file content type got through", err)
	}

}

// Tests ProcessUploadRequest, and also DeleteTrack to clean up and test that too
func TestProcessUpload(t *testing.T) {
	app := testutils.InitializeFirebaseTest()
	values := map[string]io.Reader{
		"file": mustOpen("../testdata/testIgc.igc"),
	}

	req, err := createMultipart(values, "/upload", "POST")
	req.Header.Set("uid", constant.ScraperUID)
	if err != nil {
		t.Error("Could not create multipart")
	}

	code, md, err := ProcessUploadRequest(app, req)
	if err != nil && code != http.StatusOK {
		t.Error("Could not save file, should pass", err)
	}

	md, err = TakeOwnership(app, md.TrackID, constant.TestUID)
	if err != nil && md.UID != constant.TestUID {
		t.Error("UID of returned object should be TestUID", err)
	}

	code, err = DeleteTrack(app, md.TrackID, constant.ScraperUID)
	if err == nil || code != http.StatusForbidden {
		t.Error("This deletion SHOULD FAIL! (ScraperUID should no longer own this track)")
	}

	code, err = DeleteTrack(app, md.TrackID, constant.TestUID)
	if err != nil && code != http.StatusOK {
		t.Error("Could not delete data, should delete data")
	}

}

// Shamefully nicked stackoverflow answer: https://stackoverflow.com/a/20397167/7036624
// Slightly modified to suit our needs
func createMultipart(values map[string]io.Reader, target string, method string) (*http.Request, error) {
	var b bytes.Buffer
	var req http.Request
	w := multipart.NewWriter(&b)

	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add binary
		if x, ok := r.(*os.File); ok {
			fw, err = w.CreateFormFile(key, x.Name())
			if err != nil {
				return &req, err
			}
		} else {
			// Add other fields
			fw, err = w.CreateFormField(key)
			if err != nil {
				return &req, err
			}
		}

		_, err = io.Copy(fw, r)
		if err != nil {
			return &req, err
		}
	}

	req2 := httptest.NewRequest(method, target, &b)
	req2.Header.Set("Content-Type", w.FormDataContentType())
	w.Close()

	return req2, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		logrus.Fatal("Not a file", err)
	}
	return r
}
