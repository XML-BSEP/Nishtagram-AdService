package dto

import (
	"ad_service/domain"
	"time"
)

type PostDTO struct {
	Id string `json:"id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Timestamp time.Time `json:"time" validate:"required"`
	NumOfLikes int `json:"numOfLikes" validate:"required"`
	NumOfDislikes int `json:"numOfDislikes" validate:"required"`
	NumOfComments int `json:"numOfComments" validate:"required"`
	Banned bool
	Profile domain.Profile `json:"user" validate:"required"`
	Location string `json:"location" validate:"required"`
	IsBookmarked bool `json:"isBookmarked" validate:"required"`
	IsDisliked bool `json:"isDisliked" validate:"required"`
	IsVideo bool `json:"isVideo" validate:"required"`
	IsAlbum bool `json:"isAlbum" validate:"required"`
	Media []string
	MediaType domain.MediaType
	Hashtags []string
	Mentions []string
	IsLiked bool `json:"isLiked" validate:"required"`

}

func NewPost(id string, desc string, timestamp time.Time, numOfLikes int, numOfDislikes int,
	numOfComments int, profileId string, locationName string, mentions []string, hashtags []string, media []string, mediaType string) PostDTO {

	var mentionsSlice []domain.Profile
	var hashtagsSlice []domain.HashTag

	for _, s := range mentions {
		mentionsSlice = append(mentionsSlice, domain.Profile{ID: s})
	}

	for _, s := range hashtags {
		hashtagsSlice = append(hashtagsSlice, domain.HashTag{Tag: s})
	}

	return PostDTO{
		Id: id,
		Description: desc,
		Hashtags: hashtags,
		Mentions: mentions,
		Timestamp: timestamp,
		NumOfLikes: numOfLikes,
		NumOfDislikes: numOfDislikes,
		NumOfComments: numOfComments,
		Profile: domain.Profile{ID: profileId},
		Location: locationName,
		Media: media,
	}
}
