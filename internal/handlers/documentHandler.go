package handlers

import (
	"encoding/json"
	"internal/models"
	"internal/service"

	"net/http"

	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	DocumentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocumentService: documentService,
	}
}

func (dh *DocumentHandler) UploadDocument(c echo.Context) error {
	mpform, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid multipart form."})
	}

	metaData, ok := mpform.Value["meta"]
	if !ok || len(metaData) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing meta data."})
	}

	var meta models.Meta
	if err := json.Unmarshal([]byte(metaData[0]), &meta); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to parse meta data."})
	}

	jsonData := c.FormValue("json")

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to get file."})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to open file."})
	}
	defer src.Close()

	// Определяем MIME-тип на основе первых 512 байт файла
	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read file data."})
	}

	mimeType := http.DetectContentType(fileData)
	// Если необходимо, вы можете добавить дополнительную проверку на корректность MIME-типа
	if mimeType == "application/octet-stream" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Unable to detect MIME type."})
	}

	// Создаем документ
	document := &models.Document{
		Name:     meta.Name,
		Mime:     mimeType, // Используем определенный MIME-тип
		Public:   meta.Public,
		FilePath: "",
	}

	saveDoc, err := dh.DocumentService.UploadDocument(document, fileData, meta.Grant)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to upload document."})
	}

	response := map[string]any{
		"data": map[string]any{
			"json": jsonData,
			"file": saveDoc.Name,
		},
	}
	return c.JSON(http.StatusOK, response)
}
