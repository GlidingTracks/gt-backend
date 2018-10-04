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
