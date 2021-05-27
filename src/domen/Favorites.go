package domen

import "time"

type Favorites struct {
	ID uint64
	Timestamp time.Time
	AdPost AdPost
	AdPostID uint64
	Profile Profile
	ProfileID uint64
}
