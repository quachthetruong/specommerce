package handler

import (
	"net/http"
	domain "specommerce/campaignservice/internal/core/domain/campaign"
	"specommerce/campaignservice/internal/core/ports/primary"
	"specommerce/campaignservice/pkg/sharedto/handler"

	"github.com/gin-gonic/gin"
)

type CampaignHandler interface {
	CreateCampaign(ctx *gin.Context)
}

type campaignHandler struct {
	campaignService primary.CampaignService
}

func NewCampaignHandler(campaignService primary.CampaignService) CampaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
	}
}

// CreateCampaign godoc
// @Summary Create a new campaign
// @Description Create a new marketing campaign with the provided details
// @Tags campaigns
// @Accept json
// @Produce json
// @Param campaign body CreateCampaignRequest true "Campaign information"
// @Success 200 {object} CampaignResponse "Campaign created successfully"
// @Failure 400 {object} handler.ErrorResponse "Bad request"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /v1/campaigns [post]
func (h *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var req CreateCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCampaign, err := h.campaignService.CreateCampaign(ctx, req.ToDomain())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[domain.Campaign]{
		Data: createdCampaign,
	})
}
