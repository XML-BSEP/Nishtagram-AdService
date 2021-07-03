package dto

type FollowDTO struct {
	Follower ProfileDTO `bson:"follower" json:"follower"`
	User     ProfileDTO `bson:"user" json:"user"`
}

