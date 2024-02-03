package file

import (
	"book_keeper/internal/logger"
	"book_keeper/internal/mongo"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
		logger.ErrorWithCtx(ctx, "failed creating temp file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(f.Name())

	if err = f.Close(); err != nil {
		logger.ErrorWithCtx(ctx, "error closing file", "error", err)
	}

	if err = ctx.SaveUploadedFile(file, f.Name()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileBytes, err := os.ReadFile(f.Name())
	if err != nil {
		logger.ErrorWithCtx(ctx, "error reading file to bytes", "error", err)
	}

	insertFile, err := h.DBClient.InsertFile(file.Filename, fileBytes, string(PDF))
	if err != nil {
		logger.ErrorWithCtx(ctx, "error storing file to database", "error", err)
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": insertFile})
	return

}

func (h *Handler) Delete(ctx *gin.Context) {
	fileID := ctx.Param("id")
	err := h.DBClient.DeleteFileByID(fileID)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error deleting file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})

}

func (h *Handler) GetAll(ctx *gin.Context) {
	files, err := h.DBClient.GetAllFiles(ctx)
	if err != nil {
		logger.ErrorWithCtx(ctx, "error fetching file", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": files})
}
