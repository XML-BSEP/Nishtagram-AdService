package dto

type LikeDislikePreviewDTO struct {
	PostId string
	User UserPreviewDTO
}

func NewLikeDislikePreviewDTO(postId string, dto UserPreviewDTO) LikeDislikePreviewDTO {
	return LikeDislikePreviewDTO{User: dto, PostId: postId}
}
