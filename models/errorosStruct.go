package models

import "net/http"

type ErrorResponseInterface interface {
	GetErrorResponse(code int) ErrorResponse
}

// ErrorResponse - модель ошибки
type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// NewErrorResponse - возвращает экземпляр бизнес ошибки (по тз)
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

// ErrorMap - маппинг ошибок
var ErrorMap = map[int]ErrorResponse{
	http.StatusUnauthorized:        {Code: 123, Text: "Не зарегестрирован"},
	http.StatusNotFound:            {Code: 123, Text: "Не найдено"},
	http.StatusForbidden:           {Code: 123, Text: "Отказано в доступе"},
	http.StatusInternalServerError: {Code: 123, Text: "Это грустно"},
}

// GetErrorResponse возвращает ошибку в зависимости от статуса http кода
func (er *ErrorResponse) GetErrorResponse(code int) ErrorResponse {
	if err, exists := ErrorMap[code]; exists {
		return err
	}
	return ErrorResponse{Code: 123, Text: "Что-то непонятное"}
}
