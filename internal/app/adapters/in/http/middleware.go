package httpadapter

import (
	"net/http"
	"strings"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/ports"

	"github.com/gin-gonic/gin"
)

func JWTAuth(tokenManager ports.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "invalid authorization header format",
			})
			return
		}

		tokenStr := parts[1]
		userID, role, err := tokenManager.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "invalid or expired token",
			})
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}

func RequireCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, exists := c.Get("companyID")
		if !exists || companyID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  "company scope required",
			})
			return
		}
		c.Next()
	}
}
