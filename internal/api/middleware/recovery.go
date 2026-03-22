package middleware

import (
	"log"
	"net/http"

	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("panic recovered: %v", recovered)
		response.Error(c, http.StatusInternalServerError, 1000, "internal server error")
	})
}
