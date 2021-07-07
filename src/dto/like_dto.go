package dto

type LikeDislikeDTO struct {
	UserId string `json:"userId" validate:"required"`
	PostId string `json:"postId" validate:"required"`
	PostBy string `json:"postBy" validate:"required"`
}
