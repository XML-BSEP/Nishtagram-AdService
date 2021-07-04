package dto

type CreatePostDTO struct {
	UserId UserTag `json:"user_id" validate:"required"`
	Caption string `json:"caption" validate:"required"`
	Location string `json:"location" validate:"required"`
	Hashtags []string `json:"hashtags" validate:"required"`
	Mentions []UserTag `json:"taggedUsers" validate:"required"`
	IsAlbum bool `json:"isAlbum" validate:"required"`
	IsImage bool `json:"isImage" validate:"required"`
	IsVideo bool `json:"isVideo" validate:"required"`
	Image string `json:"image"`
	Album []string `json:"album"`
	Media []string
	MediaType string
	Video string `json:"video"`
	MentionsToAdd []string

}

type UserTag struct {
	UserId string `json:"userId" validate:"required"`
	Username string `json:"username" validate:"required"`
}
