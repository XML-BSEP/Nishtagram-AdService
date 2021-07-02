package domain

import "time"

type Dislike struct {
	ID uint64
	Timestamp time.Time
	AdPostID uint64
}