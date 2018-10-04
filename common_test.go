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

func TestGetBoolFromString(t *testing.T) {
	t.Run("Return false", func(t *testing.T) {
		b := GetBoolFromString("false")
		if b {
			t.Error("Expected false, got: ", b)
		}
	})

	t.Run("Return true", func(t *testing.T) {
		b := GetBoolFromString("true")
		if !b {
			t.Error("Expected true, got: ", b)
		}
	})

	t.Run("->Empty, return false", func(t *testing.T) {
		b := GetBoolFromString("")
		if b {
			t.Error("Expected false, got: ", b)
		}
	})

	t.Run("->jibberish, return false", func(t *testing.T) {
		b := GetBoolFromString("fewfwfewfew")
		if b {
			t.Error("Expected false, got: ", b)
		}
	})
}
