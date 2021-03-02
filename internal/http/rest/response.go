package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	ErrCodeBadRequest = map[string]interface{}{
		"Code":     10004,
		"Message":  "Invalid parameter %v",
		"Message2": "The account id mismatch",
	}
)

type ErrorResponse struct {
	Error Response `json:"error"`
}

type Response struct {
	Code            int         `json:"code"`
	Message         string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
}

func WriteResponse(writer http.ResponseWriter, HTTPStatus int, data interface{}) string {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(HTTPStatus)

	byteData, _ := json.Marshal(data)
	_, err := writer.Write(byteData)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(byteData)
}
