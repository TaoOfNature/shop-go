package handler

import (
	"net/http"

	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type RecommendHandler struct {
	service *service.HomeService
}

func NewRecommendHandler(service *service.HomeService) *RecommendHandler {
	return &RecommendHandler{service: service}
}

func (h *RecommendHandler) Personalized(c *gin.Context) {
	items, err := h.service.PersonalizedRecommend(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2801, err.Error())
		return
	}
	response.Success(c, items)
}
