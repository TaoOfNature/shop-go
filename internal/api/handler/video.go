package handler

import (
	"net/http"

	"github.com/TaoOfNature/shop-go/internal/pkg/response"
	"github.com/TaoOfNature/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	service *service.VideoService
}

func NewVideoHandler(service *service.VideoService) *VideoHandler {
	return &VideoHandler{service: service}
}

func (h *VideoHandler) Recommend(c *gin.Context) {
	items, err := h.service.Recommend(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2701, err.Error())
		return
	}
	response.Success(c, items)
}
