package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	campaignHandler "specommerce/campaignservice/internal/adapters/primary/campaign/handler"
)

// Route for admin user
func adminRoutes(routerGroup *gin.RouterGroup, injector do.Injector) {
	campaign := do.MustInvoke[campaignHandler.CampaignHandler](injector)

	v1OrderGroup := routerGroup.Group("v1/campaigns")
	v1OrderGroup.POST("", campaign.CreateCampaign)
}
