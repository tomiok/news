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

func (h *Handler) GetLocationFeed(w http.ResponseWriter, r *http.Request) error {
	l1 := r.URL.Query().Get("l1")
	l2 := r.URL.Query().Get("l2")

	feed, err := h.Service.GetFeed(l1, l2)

	if err != nil {
		return err
	}

	return web.ResponseOK(w, "feed", feed)
}
