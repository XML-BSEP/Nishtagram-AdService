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
	GetAllDisposableCampaigns(ctx *gin.Context)
	GetAllMultipleCampaigns(ctx *gin.Context)
	UpdateMultipleCampaign(ctx *gin.Context)
	DeleteMultipleCampaign(ctx *gin.Context)
	DeleteDisposableCampaign(ctx *gin.Context)
}

type campaignHandler struct {
	campaignUseCase usecase.CampaignUseCase
}

func (c campaignHandler) GetAllDisposableCampaigns(ctx *gin.Context) {
	agentId, _ := middleware.ExtractUserId(ctx.Request)

	campaigns, err := c.campaignUseCase.GetAllDisposableCampaignsForAgent(ctx, agentId)
	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}
	ctx.JSON(200, campaigns)
}

func (c campaignHandler) GetAllMultipleCampaigns(ctx *gin.Context) {
	agentId, _ := middleware.ExtractUserId(ctx.Request)

	campaigns, err := c.campaignUseCase.GetAllMultipleCampaignsForAgent(ctx, agentId)
	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}
	ctx.JSON(200, campaigns)}

func (c campaignHandler) UpdateMultipleCampaign(ctx *gin.Context) {
	var campaign domain.MultipleCampaign

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignUseCase.UpdateMultipleCampaign(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignHandler) DeleteMultipleCampaign(ctx *gin.Context) {
	var campaign domain.MultipleCampaign

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignUseCase.DeleteMultipleCampaign(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignHandler) DeleteDisposableCampaign(ctx *gin.Context) {
	var campaign domain.DisposableCampaign

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId.ID, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignUseCase.DeleteDisposableCampaign(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
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
