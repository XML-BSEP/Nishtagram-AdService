package domain

import "time"

type MediaType int

const (
	PHOTO MediaType = iota
	VIDEO
)
type AdPost struct {
	ID string `json:"id"`
	Path string
	Description string
	Type MediaType
	AgentId Profile
	Timestamp time.Time
	NumOfLikes int
	NumOfDislikes int
	NumOfComments int
	Banned bool
	Link string
	HashTags []string
	Location string
}