package domen

import "time"

type Like struct {
	ID uint64
	Timestamp time.Time
	AdPostID uint64

}