package gtbackend

import (
	"github.com/Sirupsen/logrus"
	"os"
	"testing"
)

// THIS TEST FILE NEEDS TO BE UPDATED OR DELETED TO COMPLY WITH THE NEW IGCPARSER! ALL TESTS ARE t.Skip on 1st line!

func TestParse(t *testing.T) {
	t.Skip(t) // Skipping test until it is adjusted or just removed
	logrus.SetLevel(logrus.ErrorLevel)

	parser := Parser{
		Parsed: "./testdata/testIgc.igc",
	}

	md := parser.Parse()

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
	t.Skip(t) // Skipping test until it is adjusted or just removed
	file, err := os.Open("./testdata/testFileHRecords.txt")
	if err != nil {
		t.Error("Error received: ", err)
	}

	parser := Parser{
		Parsed: "",
	}

	l, err := parser.fileToLines(file)

	defer file.Close()

	if err != nil {
		t.Error("Error received: ", err)
	}

	if len(l) != 33 {
		t.Error("Wrong number of lines read")
	}
}

func TestGetHRecords(t *testing.T) {
	t.Skip(t) // Skipping test until it is adjusted or just removed
	file, err := os.Open("./testdata/testFileHRecords.txt")
	if err != nil {
		t.Error("Error received: ", err)
	}

	parser := Parser{
		Parsed: "",
	}

	l, err := parser.fileToLines(file)

	h := parser.getHRecords(l)

	if len(h) != 7 {
		t.Error("Expected len 7 got: ", len(h))
	}
}

func TestStrip(t *testing.T) {
	t.Skip(t) // Skipping test until it is adjusted or just removed
	parser := Parser{
		Parsed: "",
	}

	t.Run("Normal behaviour", func(t *testing.T) {
		testString := "Foo:bar"
		expected := "bar"

		actual := parser.strip(testString, ":")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

	t.Run("No occurrence of split rune", func(t *testing.T) {
		testString := "Foobar"
		expected := testString

		actual := parser.strip(testString, ":")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		testString := ""
		expected := testString

		actual := parser.strip(testString, ",")

		if actual != expected {
			t.Error("Expected: ", expected, ", got: ", actual)
		}
	})

}
