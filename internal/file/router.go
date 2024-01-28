package file

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes(router *gin.Engine) {
	fileGroup := router.Group("file")
	fileGroup.POST("", h.UploadFile)
	fileGroup.DELETE("/:id", h.Delete)
}
