package models

import "time"

type Block struct {
	ID        int64
	CreatedAt time.Time
	EndedAt   *time.Time
}
