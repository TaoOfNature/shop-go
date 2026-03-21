package handler

import (
	"net/http"

	"github.com/TaoOfNature/shop-go/internal/pkg/response"
	"github.com/TaoOfNature/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	service *service.HomeService
}

func NewHomeHandler(service *service.HomeService) *HomeHandler {
	return &HomeHandler{service: service}
}

func (h *HomeHandler) Banners(c *gin.Context) {
	items, err := h.service.Banners(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2601, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *HomeHandler) Categories(c *gin.Context) {
	items, err := h.service.Categories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2602, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *HomeHandler) Recommend(c *gin.Context) {
	items, err := h.service.Recommend(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2603, err.Error())
		return
	}
	response.Success(c, items)
}
