package rest

import (
	"github.com/GlidingTracks/gt-backend/constant"
	"github.com/GlidingTracks/gt-backend/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing that the start page, or well any path not defined in the backend returns properly
func TestStartPage(t *testing.T) {
	app := testutils.InitializeFirebaseTest()
	token := testutils.RetrieveFirebaseIDToken(app, constant.TestUID)

	r := CompleteRouterSetup(app)

	// Test root path
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("token", token)
	_ = testutils.TestRoute(req, r, "Root", t, http.StatusBadRequest)

	// Test random path
	req = httptest.NewRequest("GET", "/randomPath", nil)
	req.Header.Set("token", token)
	_ = testutils.TestRoute(req, r, "RandomPath", t, http.StatusNotFound)
}
