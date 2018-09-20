package rest

import "github.com/gorilla/mux"

// MuxRouteBinder contains methods for binding routes to a router.
type MuxRouteBinder interface {
	Bind(r *mux.Router)
}
