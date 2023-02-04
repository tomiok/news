package handler

import (
	"net/http"
	"news/internal/collector"
	"news/platform/web"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*collector.Service
}

func New(dbURL string) *Handler {
	service, err := collector.NewService(dbURL)

	if err != nil {
		return nil
	}

	return &Handler{Service: service}
}

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) error {
	uid := chi.URLParam(r, "articleUID")
	article, err := h.Service.GetNewsByUID(uid)

	if err != nil {
		return err
	}

	return web.ResponseOK(w, "news", article)
}
