package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/Mallbrusss/BackEntryMiddle/internal/service"
	"github.com/Mallbrusss/BackEntryMiddle/models"

	"net/http"

	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	DocumentService service.DocumentServiceInterface
	errRes          *models.ErrorResponse
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocumentService: documentService,
		errRes:          models.NewErrorResponse(),
	}
}

func (dh *DocumentHandler) UploadDocument(c echo.Context) error {
	mpform, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	metaData, ok := mpform.Value["meta"]
	if !ok || len(metaData) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	var meta models.Meta
	if err := json.Unmarshal([]byte(metaData[0]), &meta); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	jsonData := c.FormValue("json")

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusNotImplemented, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}
	defer src.Close()

	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		return c.JSON(http.StatusNotImplemented, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}

	mimeType := http.DetectContentType(fileData)
	if mimeType == "application/octet-stream" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	document := &models.Document{
		Name:     meta.Name,
		Mime:     mimeType,
		Public:   meta.Public,
		File:     meta.File,
		FilePath: "",
	}

	saveDoc, err := dh.DocumentService.UploadDocument(document, fileData, meta.Grant)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}

	response := map[string]any{
		"data": map[string]any{
			"json": jsonData,
			"file": saveDoc.Name,
		},
	}
	return c.JSON(http.StatusOK, response)
}

func (dh *DocumentHandler) DeleteDocument(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user.Login == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusUnauthorized)})
	}

	docId := c.Param("id")

	err := dh.DocumentService.DeleteDocument(docId, user.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}

	return c.JSON(http.StatusOK, echo.Map{"response": map[string]string{
		docId: "true",
	}})
}

func (dh *DocumentHandler) GetDocumentByID(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok || user.Login == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusUnauthorized)})
	}

	documentID := c.Param("id")

	document, err := dh.DocumentService.GetDocumentByID(documentID, user.Login)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}

	if document.File {
		return c.File(document.FilePath)

	}

	return c.JSON(http.StatusOK, map[string]any{"data": document})
}

func (dh *DocumentHandler) GetDocuments(c echo.Context) error {
	login := c.QueryParam("login")
	user, ok := c.Get("user").(*models.User)
	if !ok || user.Login == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusUnauthorized)})
	}

	if login == "" {
		login = user.Login
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	key := c.QueryParam("key")
	value := c.QueryParam("value")

	filter := make(map[string]any)
	if key != "" && value != "" {
		filter[key] = value
	}

	documents, err := dh.DocumentService.GetDocuments(login, filter, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusInternalServerError)})
	}
	return c.JSON(http.StatusOK, map[string]any{"data": map[string]any{"docs": documents}})
}

func (dh *DocumentHandler) AuthMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			user, err := dh.DocumentService.GetUserByToken(token)
			if err != nil {

				return c.JSON(http.StatusUnauthorized, echo.Map{"error": dh.errRes.GetErrorResponse(http.StatusUnauthorized)})
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func (dh *DocumentHandler) HeadDocument(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	c.Response().Header().Set(echo.HeaderAccept, "OK")

	return c.NoContent(http.StatusOK)
}
