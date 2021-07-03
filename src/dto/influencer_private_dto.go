package dto

type InfluencerPrivateDTO struct {
	IsInfluencer bool `bson:"is_influencer" json:"isInfluencer"`
	IsPrivate bool `bson:"is_private" json:"isPrivate"`
}
