package rest

import (
	"firebase.google.com/go"
)

// Context contains the Context of a Client session.
// Holds a instance of a firebase App.
type Context struct {
	App *firebase.App
}
