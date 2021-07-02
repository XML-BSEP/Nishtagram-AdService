package handler

import (
	"ad_service/domain"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type AdHandler interface {
	CreateAd(ctx *gin.Context)
	GetAdsByAgentId(ctx *gin.Context)
}

type adHandler struct {
	adUseCase usecase.AdPostUseCase
}

func (a adHandler) CreateAd(ctx *gin.Context) {
	var ad domain.AdPost

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&ad); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	ad.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := a.adUseCase.CreateAdPost(ctx, ad)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (a adHandler) GetAdsByAgentId(ctx *gin.Context) {
	agentId, _ := middleware.ExtractUserId(ctx.Request)

	ads, err := a.adUseCase.GetAdsByAgent(ctx, agentId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, ads)
}

func NewAdHandler(adUseCase usecase.AdPostUseCase) AdHandler {
	return &adHandler{adUseCase: adUseCase}
}