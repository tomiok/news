package web

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
}

func ResponseBadRequest(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseInternalError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseNotFound(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	r, _ := json.Marshal(HttpResponse{
		Message: msg,
		Success: false,
	})
	_, _ = w.Write(r)
}

func ResponseOK(w http.ResponseWriter, msg string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(HttpResponse{
		Message: msg,
		Success: true,
		Data:    data,
	})
	return nil
}

func ResponseCreated(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(HttpResponse{
		Message: msg,
		Success: true,
		Data:    data,
	})
}

func ResponseNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
