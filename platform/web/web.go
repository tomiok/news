package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HttpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
}

type Err struct {
	Message string
	Code    int
}

func (e Err) Error() string {
	return e.Message
}

func transform(e error) *Err {
	var webErr Err
	if errors.As(e, &webErr) {
		err := e.(Err)
		return &err
	}

	return &Err{
		Message: "errors during the request",
		Code:    http.StatusInternalServerError,
	}
}

func ReturnErr(w http.ResponseWriter, err error) {
	_err := transform(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(_err.Code)

	r, _ := json.Marshal(HttpResponse{
		Message: _err.Message,
		Success: false,
	})
	_, _ = w.Write(r)
}
