package rest

import (
	"errors"
	"firebase.google.com/go"
	"github.com/GlidingTracks/gt-backend"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

// Sets up all routing used by the backend, function public for tests wanting to do full E2E tests
func CompleteRouterSetup(app *firebase.App) (handler http.Handler) {
	ctx := Context{
		App: app,
	}

	sec := gtbackend.SecurityMiddleware{App: app}

	r := mux.NewRouter()
	r.Use(gtbackend.LogIncomingRequests)
	r.Use(sec.CheckIncomingRequests)

	userRoutes := UserHandler{
		Ctx:            ctx,
		CreateUserPage: "/createUser",
		UpdateUserPage: "/updateUser",
		DeleteUserPage: "/deleteUser",
		GetUserPage:    "/getUser",
	}

	dbRoutes := DbHandler{
		Ctx:         ctx,
		InsertTrack: "/insertTrack",
		GetTracks:   "/getTracks",
		GetTrack:    "/getTrack",
		DeleteTrack: "/deleteTrack",
	}

	userRoutes.Bind(r)
	dbRoutes.Bind(r)

	r.HandleFunc("/", startPage)

	handler = cors.AllowAll().Handler(r)

	return
}

// startPage redirects every non-existing path to url: localhost:8080/.
func startPage(w http.ResponseWriter, r *http.Request) {
	err := errors.New("page not found")
	http.Error(w, err.Error(), http.StatusBadRequest)
}