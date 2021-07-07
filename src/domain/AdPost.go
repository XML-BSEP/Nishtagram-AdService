package domain

import "time"

type MediaType int

const (
	PHOTO MediaType = iota
	VIDEO
)
type AdPost struct {
	ID string `json:"id"`
	Path string `json:"media"`
	Description string `json:"description"`
	Type MediaType `json:"type"`
	AgentId Profile
	Timestamp time.Time
	NumOfLikes int
	NumOfDislikes int
	NumOfComments int
	Banned bool
	Link string `json:"link"`
	HashTags []string `json:"hashtags"`
	Location string `json:"location"`
	CampaignId string `json:"campaignId"`
}