package handlers

import (
	"fmt"
	"os"
	"path/filepath"
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

func (h *UploadHandler) UploadFile(ctx *gin.Context) {
	if _, err := os.Stat(h.uploadPath); os.IsNotExist(err) {
		err := os.MkdirAll(h.uploadPath, 0755)
		if err != nil {
			utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to upload file")
			return
		}
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.RespondWithBadRequest(ctx, err.Error(), "Failed to upload file")
		return
	}
	defer file.Close()

	originalFileName := header.Filename
	fileExt := filepath.Ext(originalFileName)
	uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	fullPath := filepath.Join(h.uploadPath, uniqueName)

	dst, err := os.Create(fullPath)
	if err != nil {
		utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to upload file")
		return
	}
	defer dst.Close()

	if err = ctx.SaveUploadedFile(header, fullPath); err != nil {
		utils.RespondWithInternalServerError(ctx, err.Error(), "Failed to upload file")
		os.Remove(fullPath)
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", uniqueName)

	response := map[string]string{
		"file_url":  fileURL,
		"file_name": originalFileName,
		"file_size": fmt.Sprintf("%d", header.Size),
		"file_type": header.Header.Get("Content-Type"),
	}

	utils.RespondWithOK(ctx, "File uploaded successfully", response)
}
