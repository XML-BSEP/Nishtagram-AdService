package domain

import "time"

type Status int

const (
	ACCEPTED Status = iota
	REJECTED
	CREATED
)
type Type int

const (
	STORY Type = iota
	POST
)

type DisposableCampaign struct {
	ID uint64 `json:"id"`
	AgentId Profile
	ExposureDate time.Time
	Status Status
	Timestamp time.Time
	Post []AdPost
	Type Type
}
