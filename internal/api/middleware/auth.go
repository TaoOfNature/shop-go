package middleware

import (
	"net/http"
	"strings"

	"github.com/dawnstack/shop-go/internal/pkg/auth"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

const userClaimsKey = "user_claims"

func AuthRequired(tokens *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Error(c, http.StatusUnauthorized, 1001, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, 1002, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := tokens.Parse(parts[1])
		if err != nil || claims.Type != "access" {
			response.Error(c, http.StatusUnauthorized, 1003, "invalid token")
			c.Abort()
			return
		}

		c.Set(userClaimsKey, claims)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) int64 {
	claims, ok := c.Get(userClaimsKey)
	if !ok {
		return 0
	}
	userClaims, ok := claims.(*auth.Claims)
	if !ok {
		return 0
	}
	return userClaims.UserID
}
