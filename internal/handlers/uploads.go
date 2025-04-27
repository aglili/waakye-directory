package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aglili/waakye-directory/internal/utils"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadPath string
}

func NewUploadHandler(uploadPath string) *UploadHandler {
	return &UploadHandler{
		uploadPath: uploadPath,
	}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to the server (max size: 10MB)
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload (max 10MB)"
// @Success 200 {object} UploadResponse  "File uploaded successfully"
// @Failure 400 {object} BadRequestResponse "Bad request"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error"
// @Router /api/v1/uploads [post]
func (h *UploadHandler) UploadFile(ctx *gin.Context) {
	// Set maximum file size (10MB)
	const maxSize = 10 << 20 // 10 MB in bytes

	// Check if upload directory exists, create if needed
	if _, err := os.Stat(h.uploadPath); os.IsNotExist(err) {
		err := os.MkdirAll(h.uploadPath, 0755)
		if err != nil {
			utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to upload file")
			return
		}
	}

	// Limit request body size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "body size limit exceeded") {
			utils.RespondWithBadRequest(ctx, "File too large", "Maximum file size is 10MB")
			return
		}

		utils.RespondWithBadRequest(ctx, err.Error(), "Failed to get uploaded file")
		return
	}
	defer file.Close()

	if header.Size > maxSize {
		utils.RespondWithBadRequest(ctx, "File too large", "Maximum file size is 10MB")
		return
	}

	originalFileName := header.Filename
	fileExt := filepath.Ext(originalFileName)
	uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	fullPath := filepath.Join(h.uploadPath, uniqueName)

	dst, err := os.Create(fullPath)
	if err != nil {
		utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to create file")
		return
	}
	defer dst.Close()

	if err = ctx.SaveUploadedFile(header, fullPath); err != nil {
		utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to save uploaded file")
		os.Remove(fullPath)
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", uniqueName)

	response := map[string]string{
		"file_url":  fileURL, // Removed ctx.Request.Host to avoid incorrect URLs
		"file_name": originalFileName,
		"file_size": fmt.Sprintf("%d", header.Size),
		"file_type": header.Header.Get("Content-Type"),
	}

	utils.RespondWithOK(ctx, "File uploaded successfully", response)
}