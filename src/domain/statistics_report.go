package domain

import "ad_service/domain/events"

type StatisticsReport struct {
	CampaignId string `json:"campaignId"`
	Description string `json:"description"`
	AdvertisedLinks []string `json:"advertisedLinks"`
	Clicks []events.ClickEvent `json:"clicks"`
	AdvertisingCount AdvertisingCount `json:"advertisingCount"`
}
