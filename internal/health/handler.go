package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

type Response struct {
	Client string `json:"client,omitempty"`
	Status string `json:"status"`
}

const (
	UP   string = "UP"
	DOWN string = "DOWN"
)

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CheckSanity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Status: UP})
}
