package domen

import "time"

type MediaType int

const (
	PHOTO MediaType = iota
	VIDEO
)
type AdPost struct {
	ID uint64 `json:"id"`
	Path string
	Description string
	Type MediaType
	Timestamp time.Time
	NumOfLikes uint
	NumOfDislikes uint
	NumOfComments uint
	Banned bool
	Link string
	HashTags []HashTag
	Location Location
	DisposableCampaignID uint64
	MultipleCampaignID uint64
}