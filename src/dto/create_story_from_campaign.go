package dto

import (
	"ad_service/domain"
	"time"
)

type StoryType struct {
	Type string
}

type Media struct {
	Timestamp time.Time
	Path string
}
type StoryDTO struct {
	StoryId string `json:"id" validate:"required"`
	UserId string `json:"user_id" validate:"required"`
	Mentions []string `json:"mentions" validate: "required"`
	MediaPath Media `json:"storycontent" validate:"required"`
	Type string `json:"type" validate:"required"`
	Location domain.Location `json:"location" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	CloseFriends bool `json:"closefriends" validate:"required"`
	User domain.Profile `json:"user" validate:"required"`
	IsVideo bool `json:"isVideo" validate:"required"`
	Story string `json:"story"`
	NotFollowing bool `json:"notFollowing"`


}