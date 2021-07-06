package dto

import (
	"ad_service/domain"
	"time"
)

type ShowAdPost struct {
	Id string `json:"id" validate:"required"`
	Media []string `json:"images" validate:"required"`
	Type string
	UserName string
	UserSurname string
	PostBy string `json:"postby"`
	UserUsername string
	User domain.Profile `json:"user" validate:"required"`
	Location string `json:"location" validate:"required"`
	Description string `json:"description" validate:"required"`
	IsAlbum bool `json:"isAlbum" validate:"required"`
	Timestamp time.Time `json:"time" validate:"required"`
	NumOfLikes int `json:"numOfLikes" validate:"required"`
	NumOfDislikes int `json:"numOfDislikes" validate:"required"`
	NumOfComments int `json:"numOfComments" validate:"required"`
	Banned bool
	IsVideo bool `json:"isVideo" validate:"required"`
	IsBookmarked bool `json:"isBookmarked" validate:"required"`
	IsDisliked bool `json:"isDisliked" validate:"required"`
	IsLiked bool `json:"isLiked" validate:"required"`
	IsCampaign bool `json:"isCampaign"`
	Link string `json:"link"`
}
