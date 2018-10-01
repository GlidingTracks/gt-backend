package gtbackend

import (
	"bufio"
	"github.com/Sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

type MetaData struct {
	Manufacturer A
	Header       H
}

type A struct {
	ManufacturerID string
	UniqueID       string
	Additional     string
}

type H struct {
	Pilot              string
	FlightRecorderType string
	GliderType         string
	GliderID           string
	FirmwareVersion    string
	HardwareVersion    string
}

func Parse(file *os.File) (md MetaData) {
	l, _ := FileToLines(file)

	var h []string

	for i := 0; i < len(l); i++ {
		if m, _ := regexp.MatchString("^H", l[i]); m {
			h = append(h, l[i])
		}
	}

	defer file.Close()

	md.Manufacturer = parseA(l[0])
	md.Header = parseH(h)

	return
}

func parseA(record string) (man A) {
	r := record[1:]
	man.ManufacturerID = r[0:3]
	man.UniqueID = r[3:6]
	man.Additional = r[6:]

	return man
}

func parseH(hRecords []string) (header H) {
	keys := make(map[string]string)

	for i := 0; i < len(hRecords); i++ {
		keys[hRecords[i][2:5]] = Strip(hRecords[i][5:], ":")
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

func FileToLines(file *os.File) (lines []string, err error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func Strip(s string, sub string) (data string) {
	i := strings.Index(s, sub)
	if i == -1 {
		return s
	}
	return s[i+1:]
}
