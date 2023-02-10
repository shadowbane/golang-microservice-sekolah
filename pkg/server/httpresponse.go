package server

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type SimpleHttpResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type ErrorMessage struct {
	Success bool      `json:"success"`
	Data    ErrorData `json:"data"`
}

type ErrorData struct {
	Message string `json:"message"`
}

func SendAndLogError(w http.ResponseWriter, m string, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	zap.S().Errorf("%s %d: %v", m, statusCode, err.Error())

	jsonResponse := ErrorMessage{
		Success: false,
		Data: ErrorData{
			Message: err.Error(),
		},
	}

	response, _ := json.Marshal(jsonResponse)

	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}
