package domain

import "time"

type Dislike struct {
	PostId string
	Profile Profile
	Timestamp time.Time
	PostBy Profile
}

func NewDislike(postId string, profileId string, timestamp time.Time) Dislike {
	return Dislike{
		PostId: postId,
		Profile: Profile{ID: profileId},
		Timestamp: timestamp,
	}
}

