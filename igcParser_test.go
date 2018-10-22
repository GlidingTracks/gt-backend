package gtbackend

import (
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
	"testing"
)

// THIS TEST FILE NEEDS TO BE UPDATED OR DELETED TO COMPLY WITH THE NEW IGCPARSER! ALL TESTS ARE t.Skip on 1st line!

func TestParse(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)

	ll := loadTestFile("./testdata/testIgc.igc")
	if ll == "" {
		t.Error("Test file could not be loaded")
	}

	parser := Parser{
		Parsed: ll,
	}

	md, _ := parser.Parse()

	t.Run("Check header parsing", func(t *testing.T) {
		if md.Header.Pilot != "Krasimir Georgiev" {
			t.Errorf("Pilot not parsed correctly, got: %s", md.Header.Pilot)
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

	t.Run("Empty string", func(t *testing.T) {
		parser = Parser{
			Parsed: "",
		}

		_, err := parser.Parse()

		if err == nil {
			t.Error("Error expected, got none")
		}
	})
}

// Utility methods

func TestFileToLines(t *testing.T) {
	file, err := os.Open("./testdata/testFileHRecords.txt")
	if err != nil {
		t.Error("Error received: ", err)
	}

	l, err := FileToLines(file)

	defer file.Close()

	if err != nil {
		t.Error("Error received: ", err)
	}

	if len(l) != 33 {
		t.Error("Wrong number of lines read")
	}
}

func TestGetHRecords(t *testing.T) {
	ll := loadTestFile("./testdata/testFileHRecords.txt")
	if ll == "" {
		t.Error("Test file could not be loaded")
	}

	parser := Parser{
		Parsed: ll,
	}

	arr := strings.Split(ll, "\n")
	h := parser.getHRecords(arr)

	if len(h) != 7 {
		t.Error("Expected len 7 got: ", len(h))
	}
}

func TestStrip(t *testing.T) {
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

func loadTestFile(path string) (lines string) {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}

	defer f.Close()

	ll, err := FileToLines(f)
	if err != nil {
		return ""
	}

	lines = strings.Join(ll, "\n")

	return
}
