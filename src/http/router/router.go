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
	g.POST("getAdsByAgent", handler.GetAdsByAgentId)

	return router

}

