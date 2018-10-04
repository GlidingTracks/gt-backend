package gtbackend

import (
	"strconv"
	"testing"
	"time"
)

func TestGetUnixTime(t *testing.T) {
	expected := strconv.FormatInt(time.Now().Unix(), 10)
	actual := GetUnixTime()

	if expected != actual {
		t.Error("Expected: ", expected, ", got: ", actual, "\n")
	}
}
