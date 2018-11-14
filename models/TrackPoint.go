package models

// TrackPoint A TrackPoint object used to compress Track data shown to user
type TrackPoint struct {
	Time        int64   `json:"Type"`
	Latitude    float64 `json:"Latitude"`
	Longitude   float64 `json:"Longitude"`
	Valid       bool    `json:"Valid"`
	PressureAlt int     `json:"Pressure_alt"`
	GPSAlt      int     `json:"GPS_alt"`
	Accuracy    float64 `json:"Accuracy"`
	EngineRPM   float64 `json:"Engine_RPM"`
}
