package testutils

import "github.com/GlidingTracks/gt-backend/models"

// Test data for InsertTrackPoint
var InsertTrackPointTestData = models.TrackPoint{
	Time:         2,
	Latitude:     5.0,
	Longitude:    5.5,
	Valid:        false,
	Pressure_alt: 5,
	GPS_alt:      5,
	Accuracy:     5.0,
	Engine_RPM:   5.0,
}
