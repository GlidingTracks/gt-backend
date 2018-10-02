package gtbackend

import (
	"bufio"
	"github.com/Sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

// Record is the metadata information about a IGC approved FR recorded flight.
// Contains information about the flight recorder itself, as well as some selected header
// information from the flight
type Record struct {
	Manufacturer A
	Header       H
}

// A A record in IGC spec, always the first record in a IGC File. Contains information about
// FR manufacturer and ID, as well as some additional information (dunno)
type A struct {
	ManufacturerID string
	UniqueID       string
	Additional     string
}

// H record, header information from a flight.
type H struct {
	Pilot              string
	FlightRecorderType string
	GliderType         string
	GliderID           string
	FirmwareVersion    string
	HardwareVersion    string
}

// Parser holds a file, mostly done to keep as much as possible private,
// while also support testing
type Parser struct {
	File *os.File
}

// Parse - main routine for parsing a IGC-track. Returns a Record
func (parser Parser) Parse() (rec Record) {
	l, _ := parser.fileToLines(parser.File)

	h := parser.getHRecords(l)

	defer parser.File.Close()

	rec.Manufacturer = parseA(l[0])
	rec.Header = parser.parseH(h)

	return
}

func parseA(record string) (man A) {
	r := record[1:]
	man.ManufacturerID = r[0:3]
	man.UniqueID = r[3:6]
	man.Additional = r[6:]

	return man
}

func (parser Parser) parseH(hRecords []string) (header H) {
	keys := make(map[string]string)

	for i := 0; i < len(hRecords); i++ {
		keys[hRecords[i][2:5]] = parser.strip(hRecords[i][5:], ":")
	}

	for k, v := range keys {
		switch k {
		case "PLT":
			header.Pilot = v
			break
		case "FTY":
			header.FlightRecorderType = v
			break
		case "GTY":
			header.GliderType = v
			break
		case "GID":
			header.GliderID = v
			break
		case "RFW":
			header.FirmwareVersion = v
			break
		case "RHW":
			header.HardwareVersion = v
			break
		default:
			logrus.Info("Unsupported key: ", k)
		}
	}

	return
}

func (parser Parser) fileToLines(file *os.File) (lines []string, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func (parser Parser) strip(s string, sub string) (data string) {
	i := strings.Index(s, sub)
	if i == -1 {
		return s
	}
	return s[i+1:]
}

func (parser Parser) getHRecords(records []string) (h []string) {
	for i := 0; i < len(records); i++ {
		if m, _ := regexp.MatchString("^H", records[i]); m {
			h = append(h, records[i])
		}
	}
	return
}