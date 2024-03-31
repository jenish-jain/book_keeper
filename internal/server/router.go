package server

import (
	"book_keeper/internal/file"
	"book_keeper/internal/health"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-jain/logger"
)

type Handlers struct {
	HealthHandler *health.Handler
	FileHandler   *file.Handler
}

func (s *Server) InitRoutes(h Handlers) {
	router := s.routerGroups.rootRouter
	router.Use(requestIdMiddleWare)

	router.GET("/sanity", h.HealthHandler.CheckSanity)

	h.FileHandler.InitRoutes(router)
}

func requestIdMiddleWare(c *gin.Context) {
	u := uuid.New()
	c.Set("request_id", u)
	logger.DebugWithCtx(c, "request started", "method", c.Request.Method, "path", c.Request.URL.Path)
	c.Next()
	logger.DebugWithCtx(c, "request server")
}
