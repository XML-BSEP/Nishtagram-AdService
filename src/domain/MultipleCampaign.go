package domain

import "time"

type MultipleCampaign struct {
	ID string `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
	AdvertisementFrequency int `json:"frequency"`
	Post []AdPost `json:"ads"`
	AgentId Profile
	Type Type `json:"type"`
}
