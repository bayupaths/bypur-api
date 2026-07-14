package handler

import (
	"bytes"
	"io"
	"net/http"

	"bayupur-portofolio-be/internal/config"
	"bayupur-portofolio-be/internal/service"
	"bayupur-portofolio-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type StorageHandler struct {
	storageService *service.StorageService
	cfg            *config.Config
}

func NewStorageHandler(ss *service.StorageService, cfg *config.Config) *StorageHandler {
	return &StorageHandler{
		storageService: ss,
		cfg:            cfg,
	}
}

func (h *StorageHandler) CheckConnection(c *fiber.Ctx) error {
	isConnected, err := h.storageService.CheckConnection(c.Context())
	if err != nil || !isConnected {
		return response.SendError(c, "Failed to connect to R2 storage", "Storage connection failed", http.StatusServiceUnavailable)
	}
	return response.SendSuccess(c, fiber.Map{"success": true, "message": "Storage connection OK"}, "Storage connection OK", http.StatusOK)
}

func (h *StorageHandler) UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return response.SendError(c, "File not found in request (make sure form field is named 'image')", "Upload failed", http.StatusBadRequest)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.SendError(c, "Failed to read file", "Upload failed", http.StatusBadRequest)
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return response.SendError(c, "Failed to copy file data", "Upload failed", http.StatusInternalServerError)
	}
	data := buf.Bytes()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	fileInfo, err := h.storageService.UploadFile(c.Context(), fileHeader.Filename, data, contentType)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to upload file to R2", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, fileInfo, "File uploaded successfully", http.StatusCreated)
}

func (h *StorageHandler) ListFiles(c *fiber.Ctx) error {
	prefix := c.Query("prefix", "")
	var p *string
	if prefix != "" {
		p = &prefix
	}

	files, err := h.storageService.ListFiles(c.Context(), p)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to list files from R2", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, files, "Files retrieved successfully", http.StatusOK)
}

func (h *StorageHandler) GetStorageInfo(c *fiber.Ctx) error {
	info, err := h.storageService.GetStorageInfo(c.Context())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve storage statistics", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, info, "Storage info retrieved successfully", http.StatusOK)
}

func (h *StorageHandler) DeleteFile(c *fiber.Ctx) error {
	key := c.Params("*")
	if key == "" {
		var body map[string]string
		if err := c.BodyParser(&body); err == nil {
			key = body["key"]
		}
	}

	if key == "" {
		return response.SendError(c, "File key is required", "Validation failed", http.StatusBadRequest)
	}

	err := h.storageService.DeleteFile(c.Context(), key)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete file", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, fiber.Map{"key": key}, "File deleted successfully", http.StatusOK)
}
