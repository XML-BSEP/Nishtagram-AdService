package handler

import (
	"ad_service/domain/events"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type AdvertiseHandler interface {
	GetAllPostAdsForUser(ctx *gin.Context)
	GetAllStoryAdsForUser(ctx *gin.Context)
	AddClickEvent(ctx *gin.Context)
}

type advertiseHandler struct {
	advertiseUseCase usecase.AdvertiseUseCase
}

func (a advertiseHandler) AddClickEvent(ctx *gin.Context) {
	var event events.ClickEvent

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&event); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	err := a.advertiseUseCase.AddClickEvent(ctx, event)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
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
