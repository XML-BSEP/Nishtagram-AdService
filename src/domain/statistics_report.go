package domain

import "ad_service/domain/events"

type StatisticsReport struct {
	CampaignId string `json:"campaignId"`
	Description string `json:"description"`
	AdvertisedLinks []string `json:"advertisedLinks"`
	Clicks []events.ClickEvent `json:"clicks"`
	AdvertisingCount AdvertisingCount `json:"advertisingCount"`
	NumOfLikes int `json:"numOfLikes"`
	NumOfDislikes int `json:"numOfDislikes"`
	NumOfComments int `json:"numOfComments"`
	CampaignType string `xml:"campaign_type" json:"campaignType"`
	CampaignPeriod string `xml:"campaign_period" json:"campaignPeriod"`
	AdvertisementFrequency string `xml:"advertisement_frequency" json:"advertisementFrequency"`
}
