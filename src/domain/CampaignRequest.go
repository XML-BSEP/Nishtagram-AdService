package domain

type RequestStatus int

const (
	REQUESTED RequestStatus = iota
	APPROVED
	REJECTED
)

type DisposableCampaignRequest struct {
	Id string `json:"id"`
	InfluencerId string `json:"influencerId"`
	DisposableCampaign DisposableCampaign `json:"disposableCampaign"`
	RequestStatus RequestStatus
	AgentId string
}

type MultipleCampaignRequest struct {
	Id string `json:"id"`
	InfluencerId string `json:"influencerId"`
	MultipleCampaign MultipleCampaign `json:"multipleCampaign"`
	RequestStatus RequestStatus
	AgentId string
}
