package handler

import (
	"net/http"
	"specommerce/paymentservice/internal/core/ports/primary"
	"specommerce/paymentservice/pkg/sharedto/handler"

	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	GetAllPayments(ctx *gin.Context)
}
type paymentHandler struct {
	paymentService primary.PaymentService
}

func NewPaymentHandler(paymentService primary.PaymentService) PaymentHandler {
	return &paymentHandler{
		paymentService: paymentService,
	}
}

// GetAllPayments godoc
// @Summary Get all payments
// @Description Retrieve all payments from the system
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {array} PaymentResponse "List of payments"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /admin/v1/payments [get]
func (h *paymentHandler) GetAllPayments(ctx *gin.Context) {
	payments, err := h.paymentService.GetAllPayments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[[]PaymentResponse]{
		Data: ToGetAllPaymentResponse(payments),
	})
}
