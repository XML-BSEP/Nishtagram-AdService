package usecase

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/repository"
	"context"
)

type LikeUseCase interface {
	LikePost(dto dto.LikeDislikeDTO, ctx context.Context) error
	DislikePost(dto dto.LikeDislikeDTO, ctx context.Context) error
	RemoveLike(dto dto.LikeDislikeDTO, ctx context.Context) error
	RemoveDislike(dto dto.LikeDislikeDTO, ctx context.Context) error
	GetLikesForPost(postId string, ctx context.Context) ([]dto.LikeDislikePreviewDTO, error)
	GetDislikesForPost(postId string, ctx context.Context) ([]dto.LikeDislikePreviewDTO, error)
}

type likeUseCase struct {
	likeRepository repository.LikeRepo
}

func (l likeUseCase) RemoveLike(dto dto.LikeDislikeDTO, ctx context.Context) error {

	err := l.likeRepository.RemoveLike(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())

	if err != nil {
	}
	return err
}

func (l likeUseCase) LikePost(dto dto.LikeDislikeDTO, ctx context.Context) error {
	err := l.likeRepository.LikePost(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())
	if err != nil {
		return err
	}

	if l.likeRepository.SeeIfDislikeExists(dto.PostId, dto.UserId, context.Background()) {
		err = l.likeRepository.RemoveDislike(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())

		if err != nil {
			return err
		}
	}

	return nil
}

func (l likeUseCase) DislikePost(dto dto.LikeDislikeDTO, ctx context.Context) error {
	err := l.likeRepository.DislikePost(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())
	if err != nil {
		return err
	}

	if l.likeRepository.SeeIfLikeExists(dto.PostId, dto.UserId, context.Background()) {
		err = l.likeRepository.RemoveLike(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (l likeUseCase) RemoveDislike(dto dto.LikeDislikeDTO, ctx context.Context) error {
	return l.likeRepository.RemoveDislike(dto.PostId, dto.PostBy, domain.Profile{ID: dto.UserId}, context.Background())
}

func (l likeUseCase) GetLikesForPost(postId string, ctx context.Context) ([]dto.LikeDislikePreviewDTO, error) {
	likes, err := l.likeRepository.GetLikesForPost(postId, context.Background())
	if err != nil {
		return nil, err
	}
	var dislikesPreview []dto.LikeDislikePreviewDTO
	for _, d := range likes {
		dislikesPreview = append(dislikesPreview, dto.NewLikeDislikePreviewDTO(postId, dto.NewUserPreviewDTO(d.Profile.ID)))
	}
	return dislikesPreview, nil
}

func (l likeUseCase) GetDislikesForPost(postId string, ctx context.Context) ([]dto.LikeDislikePreviewDTO, error) {
	dislikes, err := l.likeRepository.GetDislikesForPost(postId, context.Background())
	if err != nil {
		return nil, err
	}
	var dislikesPreview []dto.LikeDislikePreviewDTO
	for _, d := range dislikes {
		dislikesPreview = append(dislikesPreview, dto.NewLikeDislikePreviewDTO(postId, dto.NewUserPreviewDTO(d.Profile.ID)))
	}
	return dislikesPreview, nil
}


func NewLikeUseCase(likeRepository repository.LikeRepo) LikeUseCase {
	return &likeUseCase{likeRepository: likeRepository}
}