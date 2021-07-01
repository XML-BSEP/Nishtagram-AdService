package domain

import "time"

type Like struct {
	ID string
	Timestamp time.Time
	AdPostID uint64

}