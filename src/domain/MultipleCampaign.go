package domain

import "time"

type MultipleCampaign struct {
	ID string `json:"id"`
	StartDate time.Time
	EndDate time.Time
	AdvertisementFrequency int
	Post []AdPost
	AgentId Profile
	Type Type
}
