package handler

import (
	"ad_service/domain"
	"ad_service/http/middleware"
	"ad_service/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type CampaignRequestHandler interface {
	CreateDisposableCampaignRequest(ctx *gin.Context)
	CreateMultipleCampaignRequest(ctx *gin.Context)
	ApproveDisposableCampaignRequest(ctx *gin.Context)
	ApproveMultipleCampaignRequest(ctx *gin.Context)
	RejectDisposableCampaignRequest(ctx *gin.Context)
	RejectMultipleCampaignRequest(ctx *gin.Context)
	GetAllDisposableCampaignRequests(ctx *gin.Context)
	GetAllMultipleCampaignRequests(ctx *gin.Context)
}

type campaignRequestHandler struct {
	campaignRequestUseCase usecase.CampaignRequestUseCase
}

func (c campaignRequestHandler) CreateDisposableCampaignRequest(ctx *gin.Context) {
		var campaign domain.DisposableCampaignRequest

		decoder := json.NewDecoder(ctx.Request.Body)

		if err := decoder.Decode(&campaign); err != nil {
			ctx.JSON(400, "invalid request")
			ctx.Abort()
			return
		}

		campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

		err := c.campaignRequestUseCase.CreateDisposableCampaignRequest(ctx, campaign)

		if err != nil {
			ctx.JSON(500, gin.H{"message" : "server error"})
			ctx.Abort()
			return
		}

		ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) CreateMultipleCampaignRequest(ctx *gin.Context) {
	var campaign domain.MultipleCampaignRequest

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignRequestUseCase.CreateMultipleCampaignRequest(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) ApproveDisposableCampaignRequest(ctx *gin.Context) {
	var campaign domain.DisposableCampaignRequest

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignRequestUseCase.ApproveDisposableCampaignRequest(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) ApproveMultipleCampaignRequest(ctx *gin.Context) {
	var campaign domain.MultipleCampaignRequest

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignRequestUseCase.ApproveMultipleCampaignRequest(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) RejectDisposableCampaignRequest(ctx *gin.Context) {
	var campaign domain.DisposableCampaignRequest

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignRequestUseCase.RejectDisposableCampaignRequest(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) RejectMultipleCampaignRequest(ctx *gin.Context) {
	var campaign domain.MultipleCampaignRequest

	decoder := json.NewDecoder(ctx.Request.Body)

	if err := decoder.Decode(&campaign); err != nil {
		ctx.JSON(400, "invalid request")
		ctx.Abort()
		return
	}

	campaign.AgentId, _ = middleware.ExtractUserId(ctx.Request)

	err := c.campaignRequestUseCase.RejectMultipleCampaignRequest(ctx, campaign)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, gin.H{"message" : "ok"})
}

func (c campaignRequestHandler) GetAllDisposableCampaignRequests(ctx *gin.Context) {
	userId, _ := middleware.ExtractUserId(ctx.Request)

	requests, err := c.campaignRequestUseCase.GetAllDisposableCampaignRequests(ctx, userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, requests)
}

func (c campaignRequestHandler) GetAllMultipleCampaignRequests(ctx *gin.Context) {
	userId, _ := middleware.ExtractUserId(ctx.Request)

	requests, err := c.campaignRequestUseCase.GetAllMultipleCampaignRequests(ctx, userId)

	if err != nil {
		ctx.JSON(500, gin.H{"message" : "server error"})
		ctx.Abort()
		return
	}

	ctx.JSON(200, requests)
}

func NewCampaignRequestHandler(campaignRequestUseCase usecase.CampaignRequestUseCase) CampaignRequestHandler {
	return &campaignRequestHandler{campaignRequestUseCase: campaignRequestUseCase}
}
