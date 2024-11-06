package models

import "net/http"

type ErrorResponce struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func NewErrorResponce() *ErrorResponce {
	return &ErrorResponce{}
}

var ErrorMap = map[int]ErrorResponce{
	http.StatusUnauthorized:        {Code: 123, Text: "So sad"},
	http.StatusNotFound:            {Code: 123, Text: "So sad"},
	http.StatusForbidden:           {Code: 123, Text: "So sad"},
	http.StatusInternalServerError: {Code: 123, Text: "So sad"},
}

func (er *ErrorResponce) GetErrorResponse(code int) ErrorResponce {
	if err, exists := ErrorMap[code]; exists {
		return err
	}
	return ErrorResponce{Code: 123, Text: "So sad"}
}
