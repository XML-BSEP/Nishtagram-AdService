package domain

import "time"

type Status int

type Type int

const (
	STORY Type = iota
	POST
)

type DisposableCampaign struct {
	ID string `json:"id"`
	AgentId Profile
	ExposureDate time.Time `json:"exposureDate"`
	Status Status
	Timestamp time.Time
	Post []AdPost `json:"ads"`
	Type Type `json:"type"`
}
