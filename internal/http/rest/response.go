package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	ErrCodeBadRequest = map[string]interface{}{
		"Code":    1001,
		"Message": "Invalid parameter %v",
	}
	ErrCodeURLInvalid = map[string]interface{}{
		"Code":    1002,
		"Message": "Invalid url %v",
	}
	ErrCodeDate = map[string]interface{}{
		"Code":    1003,
		"Message": "Invalid expire_date value. The expire_date must be greater than now",
	}
	ErrCodeRedis = map[string]interface{}{
		"Code":    1004,
		"Message": "Error redis",
	}
	ErrCodeUrlExp = map[string]interface{}{
		"Code":    1005,
		"Message": "url has expired",
	}
	ErrCodeNotfound = map[string]interface{}{
		"Code":    1006,
		"Message": "Data not found",
	}
)

type ErrorResponse struct {
	Error Response `json:"error"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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
