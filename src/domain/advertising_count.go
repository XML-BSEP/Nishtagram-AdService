package domain

type AdvertisingCount struct {
	CampaignId string `json:"campaignId"`
	AgentId string `json:"agentId"`
	AdvertisedCount int `json:"advertisedCount"`
}
