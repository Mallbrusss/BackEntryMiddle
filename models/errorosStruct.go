package models

import "net/http"

// ErrorResponse - модель пользователя
type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// NewErrorResponse - возвращает экземпляр кбизнес ошибки (по тз)
func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

// ErrorMap - маппинг ошибок
var ErrorMap = map[int]ErrorResponse{
	http.StatusUnauthorized:        {Code: 123, Text: "So sad"},
	http.StatusNotFound:            {Code: 123, Text: "So sad"},
	http.StatusForbidden:           {Code: 123, Text: "So sad"},
	http.StatusInternalServerError: {Code: 123, Text: "So sad"},
}

// GetErrorResponse возвращает ошибку в зависимости от статуса http кода
func (er *ErrorResponse) GetErrorResponse(code int) ErrorResponse {
	if err, exists := ErrorMap[code]; exists {
		return err
	}
	return ErrorResponse{Code: 123, Text: "So sad"}
}
