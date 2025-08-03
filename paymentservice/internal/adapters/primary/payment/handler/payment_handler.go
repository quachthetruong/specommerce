package handler

import (
	"net/http"
	"specommerce/paymentservice/internal/core/ports/primary"
	"specommerce/paymentservice/pkg/sharedto/handler"

	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	GetAllPayments(ctx *gin.Context)
	SearchPayments(ctx *gin.Context)
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

// SearchPayments godoc
// @Summary Search payments with pagination and sorting
// @Description Search payments with pagination and sorting by created_at or total_amount
// @Tags payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" minimum(1) default(1)
// @Param size query int false "Page size" minimum(1) default(10)
// @Param sort query string false "Sort by field with direction (e.g., created_at, -total_amount)"
// @Success 200 {array} PaymentResponse "Paginated payments"
// @Failure 400 {object} handler.ErrorResponse "Bad request"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /admin/v1/payments/search [get]
func (h *paymentHandler) SearchPayments(ctx *gin.Context) {
	var req SearchPaymentsRequest
	if err := handler.ParsePagination(ctx, &req.Paging); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.paymentService.SearchPayments(ctx, req.ToFilter())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
