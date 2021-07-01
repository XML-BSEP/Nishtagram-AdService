package handler

import (
	"ad_service/domain"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	CreateDisposableCampaign(ctx *gin.Context)
	CreateMultipleCampaign(ctx *gin.Context)
}

type campaignHandler struct {
	campaignUseCase usecase.CampaignUseCase
}

func (c campaignHandler) CreateDisposableCampaign(ctx *gin.Context) {
	var campaign domain.DisposableCampaign

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignUseCase.CreateDisposableCampaign(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignHandler) CreateMultipleCampaign(ctx *gin.Context) {
	var campaign domain.MultipleCampaign

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignUseCase.CreateMultipleCampaign(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func NewCampaignHandler(campaignUseCase usecase.CampaignUseCase) CampaignHandler {
	return &campaignHandler{campaignUseCase: campaignUseCase}
}
