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
// information from the flight.
type Record struct {
	Manufacturer A
	Header       H
}

// A A record in IGC spec, always the first record in a IGC File. Contains information about
// FR manufacturer and ID, as well as some additional information (dunno what it is though).
type A struct {
	ManufacturerID string
	UniqueID       string
	Additional     string
}

// H record, Metadata/header information from a FR recorded flight.
type H struct {
	Pilot              string
	FlightRecorderType string
	GliderType         string
	GliderID           string
	FirmwareVersion    string
	HardwareVersion    string
	Date               string
}

// Parser holds a file, mostly done to keep as much as possible private,
// while also support testing.
type Parser struct {
	Path string
}

// Parse - main routine for parsing a IGC-track. Returns a Record.
func (parser Parser) Parse() (rec Record, lines []string) {
	f, err := parser.openFile()
	if err != nil {
		return
	}

	lines, _ = parser.fileToLines(f)

	defer f.Close()

	if len(lines) == 0 {
		logrus.Info("No lines in file")
		return
	}

	h := parser.getHRecords(lines)

	rec.Manufacturer = parseA(lines[0])
	rec.Header = parser.parseH(h)

	return
}

// parseA processes different Axxxxxx... fields into a A object, which is returned.
func parseA(record string) (man A) {
	r := record[1:]
	man.ManufacturerID = r[0:3]
	man.UniqueID = r[3:6]
	man.Additional = r[6:]

	return man
}

// parseH processes different Hxxxxxx... fields into a H object, which is returned.
// Unsupported encountered keys is stdouted on DebugLvl info.
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
		case "DTE":
			header.Date = v[0:6]
		default:
			logrus.Info("Unsupported key: ", k)
		}
	}

	return
}

// fileToLines will traverse and return
// an array with all file lines.
func (parser Parser) fileToLines(file *os.File) (lines []string, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

// strip will take a string and remove everything before the first occurrence of the substring.
func (parser Parser) strip(s string, sub string) (data string) {
	i := strings.Index(s, sub)
	if i == -1 {
		return s
	}
	return s[i+1:]
}

// getHRecords will traverse a list of strings and return a new array with only.
// the lines starting with an capital h
func (parser Parser) getHRecords(records []string) (h []string) {
	for i := 0; i < len(records); i++ {
		if m, _ := regexp.MatchString("^H", records[i]); m {
			h = append(h, records[i])
		}
	}
	return
}

// openFile will return an open file based on the parser's Path var.
func (parser Parser) openFile() (file *os.File, err error) {
	file, err = os.Open(parser.Path)
	return
}
