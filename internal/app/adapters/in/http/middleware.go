package httpadapter

import (
	"context"
	"net/http"
	"strings"
	"time4book/internal/app/adapters/in/http/handlers"
	"time4book/internal/app/core/domain/model/company"
	"time4book/internal/app/core/domain/model/user"
	"time4book/internal/app/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		c.Set("userID", userID.String())
		c.Set("role", role)
		c.Next()
	}
}

// RequireCompanyScope scopes tenant data. Non-developers always use company_id from their profile.
// Developers bypass profile-based company and may optionally narrow with ?companyId= or X-Company-ID header.
// If a developer sends neither, list handlers receive no company scope (cross-tenant reads where supported).
func RequireCompanyScope(userRepo user.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString("role")
		if roleStr == string(user.RoleDeveloperKey) {
			if cid := strings.TrimSpace(c.Query("companyId")); cid != "" {
				if id, err := uuid.Parse(cid); err == nil {
					c.Set("companyID", id)
				}
			} else if hid := strings.TrimSpace(c.GetHeader("X-Company-ID")); hid != "" {
				if id, err := uuid.Parse(hid); err == nil {
					c.Set("companyID", id)
				}
			}
			c.Next()
			return
		}

		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "invalid user id in token",
			})
			return
		}

		u, err := userRepo.ByID(context.Background(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "user not found",
			})
			return
		}

		if u.CompanyID() == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  "user is not associated with any company",
			})
			return
		}

		c.Set("companyID", *u.CompanyID())
		c.Next()
	}
}

func RequireActiveCompany(
	userRepo user.UserRepo,
	companyRepo company.CompanyRepo,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString("role")
		if roleStr == string(user.RoleDeveloperKey) {
			c.Next()
			return
		}

		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "invalid user id in token",
			})
			return
		}

		u, err := userRepo.ByID(context.Background(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "user not found",
			})
			return
		}

		if u.CompanyID() == nil {
			c.Next()
			return
		}

		comp, err := companyRepo.ByID(context.Background(), *u.CompanyID())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.ErrorResponse{
				Status: false,
				Error:  "company not found",
			})
			return
		}

		if comp.IsBlocked() {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.ErrorResponse{
				Status: false,
				Error:  "company is blocked",
			})
			return
		}

		c.Next()
	}
}
