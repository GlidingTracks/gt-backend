package gtbackend

import (
	"strconv"
	"time"
)

// GetUnixTime returns a unix timestamp as a string for when the method got called
func GetUnixTime() (stamp string) {
	s := time.Now().Unix()
	return strconv.FormatInt(s, 10)
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
