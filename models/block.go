// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package models

import "time"

type Block struct {
	ID        int64
	CreatedAt time.Time
	EndedAt   *time.Time
}
