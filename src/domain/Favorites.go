package domain

import "time"

type Favorites struct {
	ID string
	Timestamp time.Time
	AdPost AdPost
	AdPostID string
	Profile Profile
	ProfileID string
}

