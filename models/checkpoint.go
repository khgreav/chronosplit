package models

import "time"

type Checkpoint struct {
	ID                   int64
	BlockID              int64
	ProjectID            int64
	SubjectID            int64
	PreviousCheckpointID *int64
	StartTime            time.Time
	EndTime              *time.Time
	Description          string
}
