package dto

type PostTagDTO struct {
	PostId	string `bson:"post_id" json:"post_id"`
	Hashtag	string `bson:"hashtag" json:"hashtag"`
}

type PostTagProfileDTO struct {
	PostId	string `bson:"post_id" json:"post_id"`
	Hashtag []string `bson:"hashtag" json:"hashtag"`
	ProfileId string `bson:"profile_id" json:"profile_id"`
}

type PostTagsDTO struct {
	PostProfileId	[]PostProfileId `bson:"post_profile_id" json:"post_profile_id"`
	Hashtag	string `bson:"hashtag" json:"hashtag"`
}