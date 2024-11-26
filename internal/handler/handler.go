package handler

import (
	_ "OnlineMusic/docs"
	"OnlineMusic/internal/client"
	"OnlineMusic/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	s *service.Service
	c *client.APIClient
}

func New(s *service.Service, c *client.APIClient) *Handler {
	return &Handler{
		s: s,
		c: c,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(Logging)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	songs := router.Group("/songs")
	{
		songs.GET("/", h.viewAll)
		songs.GET("/:id/lyrics", h.findText)
		songs.POST("/", h.add)
		songs.DELETE("/:id", h.delete)
		songs.PUT("/:id", h.update)
	}

	return router
}
