package rest

import (
	"fmt"
	"github.com/marni/goigc"
	"math"
	"strings"
)

// fileNameFUH - Used in debugging.
const fileNameFAH = "fileAnalyzeHandler.go"

func AnalyzeIGC(content string) (contentWithResults string, err error) {
	track, err := igc.Parse(content)
	//track, err := igc.ParseLocation(content)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
		return
	}
	initEvents, soaringEvents := analysisTimeBoundEvent(track)
	contentWithResults = saveResultAsIGCMetadata(content, initEvents, soaringEvents)
	return
}

func saveResultAsIGCMetadata(content string, results []string, soaring []string) (contentWithResults string) {
	lines := strings.Split(content, "\n")
	j := 0
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		// ignore empty lines
		if len(strings.Trim(line, " ")) < 1 {
			continue
		}
		if line[0] == 'B' {
			lines[i] = line + results[j] + soaring[j] + "\r"
			j++
		} else {
			continue
		}
	}
	contentWithResults = strings.Join(lines, "\n")
	return
}

func analysisTimeBoundEvent(track igc.Track) (initEvents []string, soaringEvents []string) {

	var peaks []int
	var troughs []int

	//find peaks and troughs
	for i := range track.Points {
		minutesRange := 1
		//determine whether this point is a peak
		if greaterThanBeforeAfter(i, track.Points, minutesRange) {
			if len(peaks) >= 1 {
				if track.Points[i].GNSSAltitude != track.Points[peaks[len(peaks)-1]].GNSSAltitude {
					peaks = append(peaks, i)
				}
			} else {
				peaks = append(peaks, i)
			}
		}

		//determine whether this point is a trough
		if lessThanBeforeAfter(i, track.Points, minutesRange) {
			if len(troughs) >= 1 {
				if track.Points[i].GNSSAltitude != track.Points[troughs[len(troughs)-1]].GNSSAltitude {
					troughs = append(troughs, i)
				}
			} else {
				troughs = append(troughs, i)
			}
		}
	}

	//fmt.Printf("peaks: %d\n", len(peaks))
	//for _, index := range (peaks) {
	//	point := track.Points[index]
	//	h, m, s := point.Time.Clock()
	//	fmt.Printf("%d:%d:%d index:%d, lat:%f, lon:%f, altitude:%d\n", h+2, m, s, index, point.Lat.Degrees(), point.Lng.Degrees(), point.GNSSAltitude)
	//}
	//
	//fmt.Printf("troughs: %d\n", len(troughs))
	//for _, index := range (troughs) {
	//	point := track.Points[index]
	//	h, m, s := point.Time.Clock()
	//	fmt.Printf("%d:%d:%d index:%d, lat:%f, lon:%f, altitude:%d\n", h+2, m, s, index, point.Lat.Degrees(), point.Lng.Degrees(), point.GNSSAltitude)
	//}

	//calculate the average speed between every peak and trough
	aveSpeed := calculateAverageSpeed(troughs, peaks, track.Points)

	//recognize the soaring events
	soaringEvents = findSoaring(track.Points)

	//use average speed to recognize init time-bound events
	for i := 0; i < len(aveSpeed); i++ {
		initEvents = append(initEvents, recogniseTimeBoundEvent(aveSpeed[i]))
	}

	// distinguish lift and thermal from flight
	_, _, initEvents = findAcceletionOfLiftThermal(aveSpeed, 1, initEvents)

	//for i,point:=range(track.Points){
	//	h, m, s := point.Time.Clock()
	//	fmt.Printf("%d:%d:%d index:%d, altitude:%d， %s\n", h+2, m, s, i,point.GNSSAltitude,initEvents[i])
	//}

	return initEvents, soaringEvents

}

// determine whether altitude of one point is greater than the altitudes of all points after minutesRange
func greaterThanBeforeAfter(index int, trackPoints []igc.Point, minutesRange int) (isTrue bool) {

	seconds := minutesRange * 60
	count := seconds / (int(trackPoints[1].Time.Unix() - trackPoints[0].Time.Unix()))
	curAltitude := trackPoints[index].GNSSAltitude

	j := 0
	if len(trackPoints) < index+count+1 {
		j = len(trackPoints)
	} else {
		j = index + count + 1
	}

	for i := index + 1; i < j; i++ {
		if curAltitude < trackPoints[i].GNSSAltitude {
			return false
		}
	}

	if -1 < index-count-1 {
		j = index - count - 1
	} else {
		j = -1
	}

	for i := index - 1; i > j; i-- {
		if curAltitude < trackPoints[i].GNSSAltitude {
			return false
		}
	}

	return true
}

// determine whether altitude of one point is less than the altitudes of all points after minutesRange
func lessThanBeforeAfter(index int, trackPoints []igc.Point, minutesRange int) (isTrue bool) {

	seconds := minutesRange * 60
	count := seconds / (int(trackPoints[1].Time.Unix() - trackPoints[0].Time.Unix()))
	curAltitude := trackPoints[index].GNSSAltitude

	j := 0
	if len(trackPoints) < index+count+1 {
		j = len(trackPoints)
	} else {
		j = index + count + 1
	}

	for i := index + 1; i < j; i++ {
		if curAltitude > trackPoints[i].GNSSAltitude {
			return false
		}
	}

	if -1 < index-count-1 {
		j = index - count - 1
	} else {
		j = -1
	}

	for i := index - 1; i > j; i-- {
		if curAltitude > trackPoints[i].GNSSAltitude {
			return false
		}
	}

	return true
}

//calculate the average speed between every peak and trough
func calculateAverageSpeed(troughs []int, peaks []int, trackPoints []igc.Point) (averageSpeed []float32) {

	temp := []int{0}
	i, j := 0, 0
	for i < len(peaks) && j < len(troughs) {
		if peaks[i] < troughs[j] {
			temp = append(temp, peaks[i])
			i++
		} else {
			temp = append(temp, troughs[j])
			j++
		}
	}
	if i == len(peaks) {
		temp = append(temp, troughs[j:]...)
	} else {
		temp = append(temp, peaks[i:]...)
	}
	temp = append(temp, len(trackPoints)-1)
	fmt.Print(temp)
	for i := 1; i < len(temp); i++ {
		altitudeChange := float32(trackPoints[temp[i]].GNSSAltitude - trackPoints[temp[i-1]].GNSSAltitude)
		timeInterval := float32(trackPoints[temp[i]].Time.Unix() - trackPoints[temp[i-1]].Time.Unix())
		speed := altitudeChange / timeInterval
		for j := temp[i-1]; j < temp[i]; j++ {
			averageSpeed = append(averageSpeed, speed)
		}
	}
	averageSpeed = append(averageSpeed, averageSpeed[i-1])
	//fmt.Print(averageSpeed)
	return
}

//distinguish lift and thermal from flight
func findAcceletionOfLiftThermal(verticalSpeed []float32, timeInterlevel int, initTimeEvents []string) (startIndex []int, endIndex []int, initEvents []string) {
	count := 0
	initEvents = initTimeEvents
	for i, s := range verticalSpeed {
		if s > 0.2 {
			count = count + 1
		} else {
			if count >= int(20/timeInterlevel) {
				//save the startIndex and endIndex
				endIndex = append(endIndex, i-1)
				startIndex = append(startIndex, i-count)
				for j := i - count; j <= i-1; j++ {
					initEvents[j] = "lift or thermal"
				}
				fmt.Println(i-count, i-1)
			}
			count = 0
		}
	}
	return
}

//initial recognise of timeBoundEvent,just distinguish flight and sinking
func recogniseTimeBoundEvent(speed float32) (timeBoundEvents string) {
	//firstly do not distinguish lift and thermal from flight
	if speed >= -1.2 {
		timeBoundEvents = "flight"
	} else if speed < -1.2 {
		if speed > -2.5 {
			timeBoundEvents = "sinking"
		} else {
			timeBoundEvents = "severe sinking"
		}
	}
	return
}

func findSoaring(trackPoints []igc.Point) (soaringEvents []string) {
	for i := 0; i < len(trackPoints); i++ {
		soaringEvents = append(soaringEvents, "")
	}
	startI := 0
	endI := 0
	for i := 0; i < len(trackPoints); i++ {
		startI = i
		for j := i + 1; j < len(trackPoints); j++ {
			if math.Abs(float64(trackPoints[j].GNSSAltitude-trackPoints[i].GNSSAltitude)) <= 1 {
				continue
			} else {
				endI = j - 1
				break

			}
		}
		if endI-startI >= 5 {
			for k := startI; k <= endI; k++ {
				soaringEvents[k] = " and soaring"
			}
		}
	}
	return
}
