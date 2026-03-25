package handler

import (
	"net/http"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type SeckillHandler struct {
	service *service.SeckillService
}

func NewSeckillHandler(service *service.SeckillService) *SeckillHandler {
	return &SeckillHandler{service: service}
}

func (h *SeckillHandler) CreateOrder(c *gin.Context) {
	var req service.CreateSeckillOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2901, err.Error())
		return
	}
	result, err := h.service.Submit(c.Request.Context(), middleware.CurrentUserID(c), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2902, err.Error())
		return
	}
	response.Success(c, result)
}
