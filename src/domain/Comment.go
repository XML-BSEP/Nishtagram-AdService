package domain

import "time"

type Comment struct {
	ID uint64
	Comment string
	Timestamp time.Time
	Profile Profile
	ProfileID uint64
	Tags []Profile
}