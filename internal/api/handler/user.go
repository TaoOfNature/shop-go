package handler

import (
	"net/http"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, err := h.service.GetProfile(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusNotFound, 2101, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2102, err.Error())
		return
	}
	if err := h.service.UpdateProfile(c.Request.Context(), middleware.CurrentUserID(c), req); err != nil {
		response.Error(c, http.StatusBadRequest, 2103, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": true})
}
