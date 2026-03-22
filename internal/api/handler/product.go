package handler

import (
	"net/http"
	"strconv"

	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Search(c *gin.Context) {
	items, err := h.service.Search(c.Request.Context(), c.Query("keyword"), c.Query("sort"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2301, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *ProductHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 2302, "invalid product id")
		return
	}
	item, err := h.service.Detail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 2303, err.Error())
		return
	}
	response.Success(c, item)
}
