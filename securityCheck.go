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
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permit X-Origin

		if checkIncomingRequests(r, sec.App) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func checkIncomingRequests(r *http.Request, app *firebase.App) bool {
	if checkIfInsecureRequest(r.RequestURI) {
		return true
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		DebugLog(InternalLog{Origin: fileNameSC, Method: "checkIncomingRequests", Err: err, Msg: "Cannot connect to Auth client"})
		return false
	}

	token, err := client.VerifyIDToken(context.Background(), r.Header.Get("token"))
	if err != nil {
		DebugLog(InternalLog{Origin: fileNameSC, Method: "checkIncomingRequests", Err: err, Msg: "Failed to authenticate user"})
		return false
	}

	// Pre-fill uid to correct UID corresponding to IDToken
	r.Header.Set("uid", token.UID)

	return true
}

// checkIfInsecureRequest - Checks if the request does not require token security check
// Using struct as it appears to be the fastest on small number of strings
func checkIfInsecureRequest(request string) bool {
	switch request {
	case
		"/createUser",
		"/updateUser",
		"/deleteUser",
		"/getUser":
		return true
	}
	return false
}
