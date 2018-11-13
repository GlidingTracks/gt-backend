package gtbackend

import (
	"context"
	"firebase.google.com/go"
	"net/http"
)

// fileNameSC filename
const fileNameSC = "securityCheck.go"

// SecurityMiddleware - capture http.Handle.
type SecurityMiddleware struct {
	App *firebase.App
}

// CheckIncomingRequests - Logs request traffic into our app.
func (sec *SecurityMiddleware) CheckIncomingRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkIncomingRequests(r, sec.App) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func checkIncomingRequests(r *http.Request, app *firebase.App) bool {
	log := DebugLogPrepareHeader(fileNameSC, "checkIncomingRequests")
	client, err := app.Auth(context.Background())
	if err != nil {
		DebugLogErrMsg(log, err, "Cannot connect to Auth client")
		return false
	}

	token, err := client.VerifyIDToken(context.Background(), r.Header.Get("token"))
	if err != nil {
		DebugLogErrMsg(log, err, "Failed to authenticate user")
		return false
	}

	// Pre-fill uid to correct UID corresponding to IDToken
	r.Header.Set("uid", token.UID)

	return true
}
