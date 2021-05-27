package domen

import "time"

type Status int

const (
	ACCEPTED Status = iota
	REJECTED
	CREATED
)

type DisposableCampaign struct {
	ID uint64 `json:"id"`
	PlacedTime time.Time
	ExposureDate time.Time
	Status Status
	Timestamp time.Time
	DisposableCampaigns []DisposableCampaign
	MultipleCampaigns []MultipleCampaign
}
