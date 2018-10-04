package gtbackend

import (
	"time"
)

// GetUnixTime returns a unix timestamp for when the method got called
func GetUnixTime() (stamp int64) {
	return time.Now().Unix()
}

// GetBoolFromString converts a string representation of bool to type bool
// default false
func GetBoolFromString(toCheck string) (b bool) {
	switch toCheck {
	case "true":
		return true
	default:
		return false
	}
}
