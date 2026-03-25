package middleware

import (
	"log/slog"
	"net/http"

	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		slog.Error("panic recovered", "error", recovered)
		response.Error(c, http.StatusInternalServerError, 1000, "internal server error")
	})
}
