package router

import (
	"ad_service/http/handler"
	"ad_service/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler handler.AppHandler) *gin.Engine {
	router := gin.Default()

	g := router.Group("/ad")

	g.Use(middleware.AuthMiddleware())

	g.POST("createAd", handler.CreateAd)
	g.POST("createDisposableCampaign", handler.CreateDisposableCampaign)
	g.POST("createMultipleCampaign", handler.CreateMultipleCampaign)
	g.GET("getAdsByAgent", handler.GetAdsByAgentId)
	g.GET("getAllDisposableCampaigns", handler.GetAllDisposableCampaigns)
	g.GET("getAllMultipleCampaigns", handler.GetAllMultipleCampaigns)
	g.POST("updateMultipleCampaign", handler.UpdateMultipleCampaign)
	g.POST("deleteMultipleCampaign", handler.DeleteMultipleCampaign)
	g.POST("deleteDisposableCampaign", handler.DeleteDisposableCampaign)
	g.POST("createDisposableCampaignRequest", handler.CreateDisposableCampaignRequest)
	g.POST("createMultipleCampaignRequest", handler.CreateMultipleCampaignRequest)
	g.POST("approveDisposableCampaignRequest", handler.ApproveDisposableCampaignRequest)
	g.POST("approveMultipleCampaignRequest", handler.ApproveMultipleCampaignRequest)
	g.POST("rejectDisposableCampaignRequest", handler.RejectDisposableCampaignRequest)
	g.POST("rejectMultipleCampaignRequest", handler.RejectMultipleCampaignRequest)
	g.GET("getAllDisposableCampaignRequests", handler.GetAllDisposableCampaignRequests)
	g.GET("getAllMultipleCampaignRequests", handler.GetAllMultipleCampaignRequests)
	g.GET("getAllPostAds", handler.GetAllPostAdsForUser)
	g.GET("getAllStoryAds", handler.GetAllStoryAdsForUser)
	return router

}

