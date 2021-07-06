package dto

import "ad_service/domain"

type CommentDTO struct {
	Comment string `json:"text" validate:"required"`
	PostId string `json:"postId" validate:"required"`
	PostBy string `json:"postBy" validate:"required"`
	CommentBy domain.Profile `json:"user" validate:"required"`
	CommentId string `json:"comment_id"`
}