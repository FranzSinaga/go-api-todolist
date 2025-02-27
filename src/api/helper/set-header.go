package api_helper

import "net/http"

func SendResponse(w http.ResponseWriter, response []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func SetResponse(data any, error bool, statusCode int, message string) map[string]any {
	response := map[string]any{
		"data":    data,
		"error":   error,
		"status":  http.StatusOK,
		"message": message,
	}
	return response
}
