package http

import (
	"chemistry-coach/internal/delivery/http/handler"
	"chemistry-coach/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Profile *handler.ProfileHandler
	Catalog *handler.CatalogHandler
	Session *handler.SessionHandler
}

func NewRouter(h Handlers, isDev bool) *gin.Engine {
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.POST("/auth/start", h.Auth.Start)
		api.GET("/goals", h.Catalog.Goals)
		api.GET("/personas", h.Catalog.Personas)

		auth := api.Group("")
		auth.Use(middleware.RequireUserID())
		{
			auth.GET("/profile", h.Profile.Get)
			auth.POST("/sessions", h.Session.Create)
			auth.GET("/sessions", h.Session.List)
			auth.POST("/sessions/:sessionId/messages", h.Session.SendMessage)
			auth.POST("/sessions/:sessionId/suggest", h.Session.Suggest)
			auth.POST("/sessions/:sessionId/finish", h.Session.Finish)
			auth.GET("/sessions/:sessionId", h.Session.Get)
		}
	}
	return r
}
