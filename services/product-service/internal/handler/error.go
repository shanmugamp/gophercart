package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	requestID string,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	response := ErrorResponse{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}

	_ = json.NewEncoder(w).Encode(response)
}
