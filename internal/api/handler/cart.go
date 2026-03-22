package handler

import (
	"net/http"
	"strconv"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2401, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *CartHandler) Add(c *gin.Context) {
	var req service.UpsertCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2402, err.Error())
		return
	}
	if err := h.service.Add(c.Request.Context(), middleware.CurrentUserID(c), req); err != nil {
		response.Error(c, http.StatusBadRequest, 2403, err.Error())
		return
	}
	response.Success(c, gin.H{"created": true})
}

func (h *CartHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2404, "invalid cart id")
		return
	}
	var req service.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2405, err.Error())
		return
	}
	if err := h.service.Update(c.Request.Context(), middleware.CurrentUserID(c), id, req); err != nil {
		response.Error(c, http.StatusBadRequest, 2406, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func (h *CartHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2407, "invalid cart id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), middleware.CurrentUserID(c), id); err != nil {
		response.Error(c, http.StatusBadRequest, 2408, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
