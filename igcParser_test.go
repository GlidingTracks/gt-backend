package gtbackend

import (
	"github.com/Sirupsen/logrus"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)

	file, err := os.Open("./testdata/testIgc.igc")
	if err != nil {
		t.Error("Could not open test file")
	}

	md := Parse(file)

	t.Run("Check header parsing", func(t *testing.T) {
		if md.Header.Pilot != "Krasimir Georgiev" {
			t.Error("Pilot not parsed correctly")
		}
	})

	t.Run("Checking A record parsing", func(t *testing.T) {
		if md.Manufacturer.ManufacturerID != "XCT" {
			t.Error("Expected XCT, got: ", md.Manufacturer.ManufacturerID)
		}

		if md.Manufacturer.UniqueID != "5e1" {
			t.Error("Expected 5e1, got: ", md.Manufacturer.UniqueID)
		}
	})

}

// Utility methods

func TestFileToLines(t *testing.T) {
	file, err := os.Open("./testdata/text.txt")
	if err != nil {
		t.Error("Error received: ", err)
	}

	l, err := FileToLines(file)

	defer file.Close()

	if err != nil {
		t.Error("Error received: ", err)
	}

	if len(l) != 3 {
		t.Error("Wrong number of lines read")
	}
}

func TestStrip(t *testing.T) {
	t.Run("Normal behaviour", func(t *testing.T) {
		testString := "Foo:bar"
		expected := "bar"

		actual := Strip(testString, ":")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

	t.Run("No occurrence of split rune", func(t *testing.T) {
		testString := "Foobar"
		expected := testString

		actual := Strip(testString, ":")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		testString := ""
		expected := testString

		actual := Strip(testString, ",")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

}
