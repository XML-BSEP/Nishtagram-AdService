package usecase

import (
	"ad_service/dto"
	"ad_service/repository"
	"context"
)

type CommentUseCase interface {
	AddComment(comment dto.CommentDTO, ctx context.Context) error
	DeleteComment(comment dto.CommentDTO, ctx context.Context) error
	GetAllCommentsByPost(postId string, ctx context.Context) ([]dto.CommentDTO, error)
}

type commentUseCase struct {
	commentRepository repository.CommentRepo
}

func (c commentUseCase) GetAllCommentsByPost(postId string, ctx context.Context) ([]dto.CommentDTO, error) {
	comments, err := c.commentRepository.GetComments(postId, context.Background())
	if err != nil {
	}

	//for i, comment := range comments {
	//	post,err := c.postRepository.get
	//	profile, err := gateway.GetUser(context.Background(), post.Profile.Id)
	//	if err != nil {
	//		p.logger.Logger.Errorf("error while getting user info for %v, error: %v\n", post.Profile.Id, err)
	//	}
	//	 comments[i].CommentBy.ProfilePhoto = profile.ProfilePhoto
	//	 comments[i].CommentBy.Username = profile.Username
	//
	//}
	return comments, err
}

func (c commentUseCase) AddComment(comment dto.CommentDTO, ctx context.Context) error {
	err := c.commentRepository.CommentPost(comment, context.Background())
	return err
}

func (c commentUseCase) DeleteComment(comment dto.CommentDTO, ctx context.Context) error {
	err := c.commentRepository.DeleteComment(comment, context.Background())
	if err != nil {
	}
	return err
}


func NewCommentUseCase(commentRepository repository.CommentRepo) CommentUseCase {
	return &commentUseCase{
		commentRepository: commentRepository,
	}
}