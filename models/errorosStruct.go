package models

import "net/http"

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func NewErrorResponse() *ErrorResponse {
	return &ErrorResponse{}
}

var ErrorMap = map[int]ErrorResponse{
	http.StatusUnauthorized:        {Code: 123, Text: "So sad"},
	http.StatusNotFound:            {Code: 123, Text: "So sad"},
	http.StatusForbidden:           {Code: 123, Text: "So sad"},
	http.StatusInternalServerError: {Code: 123, Text: "So sad"},
}

func (er *ErrorResponse) GetErrorResponse(code int) ErrorResponse {
	if err, exists := ErrorMap[code]; exists {
		return err
	}
	return ErrorResponse{Code: 123, Text: "So sad"}
}
