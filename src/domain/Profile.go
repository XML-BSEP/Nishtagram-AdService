package domain

type Profile struct {
	ID string `json:"profileId"`
	ProfileId string `json:"id"`
	Username string `json:"username" validate:"required"`
	ProfilePhoto string `json:"profilePhoto" validate:"required"`
}