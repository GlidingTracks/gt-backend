package models

type IgcRecordMetadata struct {
	DummyString string
}

// FilePayload - Client side payload for inserting a track.
type IgcMetadata struct {
	Privacy  bool
	Time string
	UID string
	Record IgcRecordMetadata
}

