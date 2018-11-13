package models

// TrackPoint A TrackPoint object used to compress Track data shown to user
type TrackPoint struct {
	Time         int64
	Latitude     float64
	Longitude    float64
	Valid        bool
	Pressure_alt int
	GPS_alt      int
	Accuracy     float64
	Engine_RPM   float64
}
