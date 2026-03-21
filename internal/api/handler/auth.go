package handler

import (
	"net/http"

	"github.com/TaoOfNature/shop-go/internal/pkg/response"
	"github.com/TaoOfNature/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2001, err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2002, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2003, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 2004, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req service.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 2005, err.Error())
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 2006, err.Error())
		return
	}
	response.Success(c, result)
}
