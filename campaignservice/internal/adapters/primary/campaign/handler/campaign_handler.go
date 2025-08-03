package handler

import (
	"net/http"
	"specommerce/campaignservice/config"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/pkg/sharedto/handler"

	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	CreateIphoneCampaign(ctx *gin.Context)
	GetIphoneWinner(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService primary.CampaignService
	config          config.AppConfig
}

func NewCampaignHandler(campaignService primary.CampaignService, config config.AppConfig) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
		config:          config,
	}
}

// CreateCampaign godoc
// @Summary Create a new campaign
// @Description Create a new marketing campaign with the provided details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body CreateCampaignRequest true "Campaign information"
// @Success 200 {object} campaign.Campaign "Campaign created successfully"
// @Failure 400 {object} handler.ErrorResponse "Bad request"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /v1/campaigns/iphones [post]
func (h *campaignHandler) CreateIphoneCampaign(ctx *gin.Context) {
	var req CreateIphoneCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCampaign, err := h.campaignService.CreateCampaign(ctx, req.ToDomain(h.config.IphoneCampaign))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[domain.Campaign]{
		Data: createdCampaign,
	})
}

// GetIphoneWinner godoc
// @Summary Get iPhone campaign winners
// @Description Get the list of winners for an iPhone campaign
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaignType path string true "Campaign type"
// @Success 200 {array} campaign.IphoneWinner "Winners retrieved successfully"
// @Failure 400 {object} handler.ErrorResponse "Bad request"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /admin/v1/campaigns/iphones/winners [get]
func (h *campaignHandler) GetIphoneWinner(ctx *gin.Context) {
	winners, err := h.campaignService.GetIphoneWinner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[[]domain.IphoneWinner]{
		Data: winners,
	})
}
