package gtbackend

import (
	"github.com/GlidingTracks/gt-backend/testutils"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityMiddleware_CheckIncomingRequests(t *testing.T) {
	app, token := testutils.RetrieveFirebaseIDTokenCredFile(false)

	sec := SecurityMiddleware{App: app}

	server := mux.NewRouter()
	server.Use(sec.CheckIncomingRequests)

	req, err := http.NewRequest("GET", "/getTracks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("token", token)

	rr := httptest.NewRecorder()
	server.HandleFunc("/getTracks", MockHandler)

	server.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Error("Wrong code returned")
	}
}
