package handler

import (
	"net/http"
	"specommerce/orderservice/internal/core/ports/primary"
	"specommerce/orderservice/pkg/sharedto/handler"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	CreateOrder(ctx *gin.Context)
	GetAllOrders(ctx *gin.Context)
}
type orderHandler struct {
	orderService primary.OrderService
}

func NewOrderHandler(orderService primary.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the provided details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order information"
// @Success 200 {object} OrderResponse "Order created successfully"
// @Failure 400 {object} handler.ErrorResponse "Bad request"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /v1/orders [post]
func (h *orderHandler) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdOrder, err := h.orderService.CreateOrder(ctx, req.ToDomain())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[OrderResponse]{
		Data: ToCreateOrderResponse(createdOrder),
	})
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Retrieve all orders from the system
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} OrderResponse "List of orders"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /admin/v1/orders [get]
func (h *orderHandler) GetAllOrders(ctx *gin.Context) {
	orders, err := h.orderService.GetAllOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, handler.BaseResponse[[]OrderResponse]{
		Data: ToGetAllOrderResponse(orders),
	})
}
