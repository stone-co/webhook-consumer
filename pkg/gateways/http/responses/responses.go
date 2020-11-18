package responses

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func Send(w http.ResponseWriter, response interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(Error{
		Message: message,
	})
}
