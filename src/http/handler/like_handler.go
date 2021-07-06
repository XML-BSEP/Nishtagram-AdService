package handler

import (
	"ad_service/dto"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type LikeHandler interface {
	LikePost(context *gin.Context)
	DislikePost(context *gin.Context)
	RemoveLike(context *gin.Context)
	RemoveDislike(ctx *gin.Context)
}

type likeHandler struct {
	likeUseCase usecase.LikeUseCase
}

func (l likeHandler) DislikePost(context *gin.Context) {
	var dislikeDTO dto.LikeDislikeDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&dislikeDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}
	dislikeDTO.UserId, _ = middleware.ExtractUserId(context.Request)


	err := l.likeUseCase.DislikePost(dislikeDTO, context)

	if err != nil {
		context.JSON(500, "server error")
		context.Abort()
		return
	}

	context.JSON(200, "ok")
}

func (l likeHandler) RemoveLike(context *gin.Context) {
	var likeDTO dto.LikeDislikeDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&likeDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}

	likeDTO.UserId, _ = middleware.ExtractUserId(context.Request)


	err := l.likeUseCase.RemoveLike(likeDTO, context)

	if err != nil {
		context.JSON(500, "server error")
		context.Abort()
		return
	}

	context.JSON(200, "ok")
}

func (l likeHandler) RemoveDislike(context *gin.Context) {
	var dislikeDTO dto.LikeDislikeDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&dislikeDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}



	dislikeDTO.UserId, _ = middleware.ExtractUserId(context.Request)
	err := l.likeUseCase.RemoveDislike(dislikeDTO, context)

	if err != nil {
		context.JSON(500, "server error")
		context.Abort()
		return
	}

	context.JSON(200, "ok")
}

func (l likeHandler) LikePost(context *gin.Context) {
	var likeDTO dto.LikeDislikeDTO

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&likeDTO); err != nil {
		context.JSON(400, "invalid request")
		context.Abort()
		return
	}

	likeDTO.UserId, _ = middleware.ExtractUserId(context.Request)


	err := l.likeUseCase.LikePost(likeDTO, context)

	if err != nil {
		context.JSON(500, "server error")
		context.Abort()
		return
	}

	context.JSON(200, "ok")
}

func NewLikeHandler(likeUseCase usecase.LikeUseCase) LikeHandler {
	return &likeHandler{likeUseCase: likeUseCase}
}
