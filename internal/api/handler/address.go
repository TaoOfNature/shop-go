package handler

import (
	"net/http"
	"strconv"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	service *service.AddressService
}

func NewAddressHandler(service *service.AddressService) *AddressHandler {
	return &AddressHandler{service: service}
}

func (h *AddressHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2201, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *AddressHandler) Create(c *gin.Context) {
	var req service.UpsertAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2202, err.Error())
		return
	}
	if err := h.service.Create(c.Request.Context(), middleware.CurrentUserID(c), req); err != nil {
		response.Error(c, http.StatusBadRequest, 2203, err.Error())
		return
	}
	response.Success(c, gin.H{"created": true})
}

func (h *AddressHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2204, "invalid address id")
		return
	}
	var req service.UpsertAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2205, err.Error())
		return
	}
	if err := h.service.Update(c.Request.Context(), middleware.CurrentUserID(c), id, req); err != nil {
		response.Error(c, http.StatusBadRequest, 2206, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2207, "invalid address id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), middleware.CurrentUserID(c), id); err != nil {
		response.Error(c, http.StatusBadRequest, 2208, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
