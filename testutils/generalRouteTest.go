package testutils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function used to send a prepared request to the handler in tests
func TestRoute(req *http.Request, r http.Handler, methodName string, t *testing.T, expectedReturnCode int) (parsedBody []byte) {
	// Run route
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Result().StatusCode != expectedReturnCode {
		t.Error("Failed " + methodName)
	}

	// Extract data from response body
	parsedBody, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Error("Failed reading body of " + methodName)
	}

	return
}
