package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Mallbrusss/BackEntryMiddle/internal/service"
	"github.com/Mallbrusss/BackEntryMiddle/models"

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

	fmt.Print("-----------------\n", metaData, "--------------\n")

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

	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read file data."})
	}

	mimeType := http.DetectContentType(fileData)
	if mimeType == "application/octet-stream" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Unable to detect MIME type."})
	}

	// Создаем документ
	document := &models.Document{
		Name:     meta.Name,
		Mime:     mimeType,
		Public:   meta.Public,
		File:     meta.File,
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

func (dh *DocumentHandler) DeleteDocument(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if user.Login == "" || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user unauthorized"})
	}

	docId := c.Param("id")

	err := dh.DocumentService.DeleteDocument(docId, user.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при удалении документа"})
	}

	return c.JSON(http.StatusOK, echo.Map{"response": docId})
}

func (dh *DocumentHandler) GetDocumentByID(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if user.Login == "" || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user unauthorized"})
	}

	documentID := c.Param("id")

	document, err := dh.DocumentService.GetDocumentByID(documentID, user.Login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if document.File {
		return c.File(document.FilePath)
	}

	return c.JSON(http.StatusOK, map[string]any{"data": document})
}

// Получение списка документов
func (dh *DocumentHandler) GetDocuments(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	if user.Login == "" || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user unauthorized"})
	}

	login := c.QueryParam("login")
	if login == "" {
		login = user.Login
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 0
	}
	key := c.QueryParam("key")
	value := c.QueryParam("value")

	filter := make(map[string]any)
	if key != "" && value != "" {
		filter[key] = value
	}

	documents, err := dh.DocumentService.GetDocuments(login, filter, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"data": map[string]any{"docs": documents}})
}

func (dh *DocumentHandler) AuthMiddleWare() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			user, err := dh.DocumentService.GetUserByToken(token)
			if err != nil {
				errResp := models.ErrorResponce{
					Code: http.StatusUnauthorized,
					Text: "Invalid token",
				}
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": errResp,
				})
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

// func (dh *DocumentHandler) HeadDocument(c echo.Context) error {
// 	ctx := context.Background()
// 	id := c.Param("id")
// 	cacheKey := fmt.Sprintf("document:%s", id)

// 	// Попытка получить данные из кэша, чтобы узнать размер и другие заголовки
// 	cachedData, err := h.RedisClient.Get(ctx, cacheKey).Result()
// 	if err != nil && err != redis.Nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Ошибка при получении данных из кэша"})
// 	}

// 	if cachedData == "" {
// 		// Данных нет в кэше, извлекаем документ из базы данных для проверки его существования
// 		document, err := h.DocumentService.GetDocumentByID(id)
// 		if err != nil {
// 			return c.JSON(http.StatusNotFound, echo.Map{"error": "Документ не найден"})
// 		}
// 		// Устанавливаем заголовки на основе данных из базы
// 		c.Response().Header().Set(echo.HeaderContentType, document.Mime)
// 		c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(document.FilePath))) // например, или другой метаданные
// 	} else {
// 		// Если данные есть в кэше, возвращаем только заголовки из кэша
// 		c.Response().Header().Set(echo.HeaderContentType, "application/json") // тип JSON данных
// 		c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(cachedData)))
// 	}

// 	// HEAD-запрос не возвращает тела
// 	return c.NoContent(http.StatusOK)
// }
