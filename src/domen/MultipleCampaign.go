package domen

import "time"

type MultipleCampaign struct {
	ID uint64 `json:"id"`
	StartDate time.Time
	EndDate time.Time
	AdvertisementFrequency 	uint
	ExposureDates []time.Time
}
