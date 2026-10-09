package httpadapter

import (
	"net/http"
	"strings"
	"time4book/internal/app/adapters/in/http/handlers"
	outports "time4book/internal/app/core/ports/out"
	userout "time4book/internal/app/core/ports/out/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func JWTAuth(tokenManager outports.TokenManager) gin.HandlerFunc {
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

		c.Set("userID", userID.String())
		c.Set("role", role)
		c.Next()
	}
}

func RequireCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, exists := c.Get("companyID")
		if !exists || companyID == nil {
			companyIDStr := c.GetHeader("x-company-id")
			if companyIDStr == "" {
				companyIDStr = c.Query("companyId")
			}

			if companyIDStr == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{
					Status: false,
					Error:  "company scope required",
				})
				return
			}

			parsedCompanyID, err := uuid.Parse(companyIDStr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, handlers.ErrorResponse{
					Status: false,
					Error:  "invalid company id",
				})
				return
			}

			c.Set("companyID", parsedCompanyID.String())
		}
		c.Next()
	}
}

func ResolveCompany(userRepo userout.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.GetString("userID"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "invalid user"})
			return
		}

		usr, err := userRepo.Retrieve(c.Request.Context(), userID)
		if err != nil || !usr.IsActive() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{Status: false, Error: "user not found or inactive"})
			return
		}

		selected := c.GetHeader("X-Company-ID")
		if selected == "" {
			selected = c.Query("companyId")
		}
		if selected != "" {
			companyID, err := uuid.Parse(selected)
			if err != nil || companyID == uuid.Nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, handlers.ErrorResponse{Status: false, Error: "invalid company id"})
				return
			}
			if !usr.Role().IsDeveloper() && companyID != usr.CompanyID() {
				c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: "company mismatch"})
				return
			}
			c.Set("companyID", companyID.String())
		} else if !usr.Role().IsDeveloper() {
			if usr.CompanyID() == uuid.Nil {
				c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{Status: false, Error: "company scope required"})
				return
			}
			c.Set("companyID", usr.CompanyID().String())
		}
		c.Next()
	}
}
