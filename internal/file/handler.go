package file

import (
	"book_keeper/internal/logger"
	"book_keeper/internal/mongo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Handler struct {
	DBClient mongo.Client
}

func NewHandler(dbClient mongo.Client) *Handler {
	return &Handler{
		DBClient: dbClient,
	}
}

func (h *Handler) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f, err := os.CreateTemp("uploads", "temp_*.pdf")
	if err != nil {
		logger.ErrorWithCtx(ctx, "failed creating temp file %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(f.Name())

	if err = f.Close(); err != nil {
		logger.ErrorWithCtx(ctx, "error closing file", fmt.Sprintf("error : %+v", err.Error()))
	}

	if err = ctx.SaveUploadedFile(file, f.Name()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileBytes, err := os.ReadFile(f.Name())
	if err != nil {
		logger.ErrorWithCtx(ctx, "error reading file to bytes", fmt.Sprintf("error : %+v", err.Error()))
	}

	insertFile, err := h.DBClient.InsertFile(ctx, file.Filename, fileBytes, string(PDF))
	if err != nil {
		logger.ErrorWithCtx(ctx, "error storing file to mongo", fmt.Sprintf("error : %+v", err.Error()))
	}
	ctx.JSON(http.StatusOK, gin.H{"data": insertFile})
	return

}

func (h *Handler) Delete(ctx *gin.Context) {
	fileID := ctx.Param("id")
	err := h.DBClient.DeleteFileByID(ctx, fileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})

}

func (h *Handler) GetAll(ctx *gin.Context) {
	files, err := h.DBClient.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": files})
}
