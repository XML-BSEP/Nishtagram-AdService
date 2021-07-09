package events

type ClickEvent struct {
	CampaignId string `json:"campaignId"`
	InfluencerId string `json:"influencerId"`
	Clicks int `json:"clicks"`
	InfluencerUsername string `json:"influencerUsername"`
}
