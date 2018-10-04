package models

import "github.com/GlidingTracks/gt-backend"

// IgcMetadata - Contains all metadata from a IGC file as well as some
// additional data.
type IgcMetadata struct {
	Privacy bool
	Time    string
	UID     string
	Record  gtbackend.Record
	TrackID string
}
