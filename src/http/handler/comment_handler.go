package handler

import (
	"ad_service/dto"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type CommentHandler interface {
	AddComment(context *gin.Context)
	DeleteComment(context *gin.Context)
	GetComments(ctx *gin.Context)
}

type commentHandler struct {
	commentUseCase usecase.CommentUseCase
}

func (c commentHandler) GetComments(ctx *gin.Context) {
	var postDTO dto.PostDTO

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&postDTO); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	comments, err := c.commentUseCase.GetAllCommentsByPost(postDTO.Id, context.Background())

	if err != nil {
		ctx.JSON(500, "server error")
		ctx.Abort()
		return
	}

	ctx.JSON(200, comments)
}

func (c commentHandler) DeleteComment(context *gin.Context) {
	var commentDTO dto.CommentDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&commentDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}

	err := c.commentUseCase.DeleteComment(commentDTO, context)

	if err != nil {
		context.JSON(500, "server error")
		context.Abort()
		return
	}

	context.JSON(200, "ok")
}

func (c commentHandler) AddComment(context *gin.Context) {
	var commentDTO dto.CommentDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&commentDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}




	commentDTO.CommentBy.ID, _ = middleware.ExtractUserId(context.Request)
	err := c.commentUseCase.AddComment(commentDTO, context)

	if err != nil {
		context.JSON(500, "server error")
	}

	context.JSON(200, "ok")
}

func NewCommentHandler(commentUseCase usecase.CommentUseCase) CommentHandler {
	return &commentHandler{commentUseCase: commentUseCase}
}