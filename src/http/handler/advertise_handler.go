package handler

import (
	"ad_service/http/middleware"
	"ad_service/usecase"
	"github.com/gin-gonic/gin"
)

type AdvertiseHandler interface {
	GetAllPostAdsForUser(ctx *gin.Context)
	GetAllStoryAdsForUser(ctx *gin.Context)
}

type advertiseHandler struct {
	advertiseUseCase usecase.AdvertiseUseCase
}

func (a advertiseHandler) GetAllPostAdsForUser(ctx *gin.Context) {
	userId, _ := middleware.ExtractUserId(ctx.Request)

	ads, err := a.advertiseUseCase.GetAllPostAdsForUser(ctx, userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, ads)
}

func (a advertiseHandler) GetAllStoryAdsForUser(ctx *gin.Context) {
	userId, _ := middleware.ExtractUserId(ctx.Request)

	ads, err := a.advertiseUseCase.GetAllStoryAdsForUser(ctx, userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, ads)
}

func NewAdvertiseHandler(advertiseUseCase usecase.AdvertiseUseCase) AdvertiseHandler {
	return &advertiseHandler{advertiseUseCase: advertiseUseCase}
}
