package handler

import (
	"net/http"
	"strconv"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2501, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *OrderHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2502, "invalid order id")
		return
	}
	order, err := h.service.Detail(c.Request.Context(), middleware.CurrentUserID(c), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 2503, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2504, err.Error())
		return
	}
	order, err := h.service.Create(c.Request.Context(), middleware.CurrentUserID(c), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2505, err.Error())
		return
	}
	response.Success(c, order)
}
