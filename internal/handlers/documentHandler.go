package handlers

import (
	"encoding/json"
	"fmt"
	"internal/models"
	"internal/service"

	"net/http"

	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	DocumentService *service.DocumentService
}

func MewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocumentService: documentService,
	}
}

func (dh *DocumentHandler) UploadDocument(c echo.Context) error {
	var meta models.Meta
	mpform, err := c.MultipartForm()
	//FIXME: Переделать под мультиформ
	if err != nil{
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}
	fmt.Println(mpform)

	metaJSON := c.FormValue("meta")
	if err := json.Unmarshal([]byte(metaJSON), &meta); err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": errResp,
		})
	}

	jsonData := c.FormValue("json")

	file, err := c.FormFile("file")
	if err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": errResp,
		})
	}

	src, err := file.Open()
	if err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": errResp})
	}
	defer src.Close()

	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": errResp})
	}

	document := &models.Document{
		Name:     meta.Name,
		Mime:     meta.Mime,
		Public:   meta.Public,
		FilePath: "", // Путь будет добавлен в DocumentService после сохранения файла
		Grant:    meta.Grant,
	}

	saveDoc, err := dh.DocumentService.UploadDocument(document, fileData, meta.Grant)
	if err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": errResp})
	}
	responce := map[string]any{
		"data": map[string]any{
			"json": jsonData,
			"file": saveDoc.Name,
		},
	}
	return c.JSON(http.StatusOK, responce)
}
