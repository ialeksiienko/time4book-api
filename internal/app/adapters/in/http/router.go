package httpadapter

import (
	"net/http"
	"time"
	"time4book/internal/app/adapters/in/http/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(h *Handler, authMw gin.HandlerFunc, companyMw gin.HandlerFunc, activeCompanyMw gin.HandlerFunc) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", h.AuthHandler.Register)
			authGroup.POST("/login", h.AuthHandler.Login)
			authGroup.POST("/refresh", h.AuthHandler.Refresh)
			authGroup.GET("/me", authMw, activeCompanyMw, h.AuthHandler.Me)
			authGroup.POST("/logout", authMw, h.AuthHandler.Logout)
		}

		userGroup := api.Group("/users", authMw, activeCompanyMw, companyMw)
		{
			userGroup.GET("", h.UserHandler.List)
			userGroup.POST("", h.UserHandler.Create)
			userGroup.PUT("/:id", h.UserHandler.Update)
			userGroup.DELETE("/:id", h.UserHandler.Deactivate)
		}

		companyGroup := api.Group("/companies", authMw, activeCompanyMw)
		{
			companyGroup.POST("", h.CompanyHandler.Create)
			companyGroup.GET("", h.CompanyHandler.List)
			companyGroup.GET("/:id", h.CompanyHandler.GetByID)
			companyGroup.PUT("/:id", h.CompanyHandler.Update)
			companyGroup.DELETE("/:id", h.CompanyHandler.Delete)
			companyGroup.POST("/:id/block", h.CompanyHandler.Block)
			companyGroup.POST("/:id/unblock", h.CompanyHandler.Unblock)
		}

		resourceGroup := api.Group("/resources", authMw, activeCompanyMw, companyMw)
		{
			resourceGroup.GET("", h.ResourceHandler.List)
			resourceGroup.POST("", h.ResourceHandler.Create)
			resourceGroup.GET("/:id", h.ResourceHandler.GetByID)
			resourceGroup.PUT("/:id", h.ResourceHandler.Update)
			resourceGroup.DELETE("/:id", h.ResourceHandler.Delete)
			resourceGroup.POST("/:id/service", h.ResourceHandler.Service)
			resourceGroup.POST("/:id/restore", h.ResourceHandler.Restore)
		}

		crtGroup := api.Group("/company-resource-types", authMw, activeCompanyMw, companyMw)
		{
			crtGroup.GET("", h.CompanyResourceTypeHandler.List)
			crtGroup.POST("", h.CompanyResourceTypeHandler.Create)
			crtGroup.PUT("/:id", h.CompanyResourceTypeHandler.Update)
			crtGroup.DELETE("/:id", h.CompanyResourceTypeHandler.Delete)
		}

		reservationGroup := api.Group("/reservations", authMw, activeCompanyMw, companyMw)
		{
			reservationGroup.GET("", h.ReservationHandler.List)
			reservationGroup.POST("", h.ReservationHandler.Create)
			reservationGroup.GET("/my", h.ReservationHandler.ListMy)
			reservationGroup.GET("/resource/:id", h.ReservationHandler.ListByResource)
			reservationGroup.POST("/:id/cancel", h.ReservationHandler.Cancel)
		}
	}

	api.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, handlers.SuccessResponse{
			Status:  true,
			Message: "ok",
		})
	})

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
				Status: false,
				Error:  err,
			})
		} else {
			c.JSON(http.StatusInternalServerError, handlers.ErrorResponse{
				Status: false,
				Error:  "internal server error",
			})
		}
		c.Abort()
	}))

	return r
}
