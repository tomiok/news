package handler

import (
	"net/http"
	"news/internal/collector"
	"news/platform/web"

	"github.com/go-chi/chi/v5"
)

// Handler will carry the services logic to the web layer.
type Handler struct {
	*collector.Service
}

// New returns a *Handler if the service is created OK, otherwise an error.
func New(dbURL string) (*Handler, error) {
	service, err := collector.NewService(dbURL)

	if err != nil {
		return nil, err
	}

	return &Handler{Service: service}, nil
}

// GetNews is a web handler that will return a news by the UID, provided in the path param.
// I.E. /news/i2W6uBcVvyDPiSCn4iMo8G
// An error is returned if the ID is not found.
func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) error {
	uid := chi.URLParam(r, "articleUID")

	if uid == "" {
		return web.Err{
			Message: "uid should be present",
			Code:    http.StatusBadRequest,
		}
	}

	article, err := h.Service.GetNewsByUID(uid)

	if err != nil {
		return err
	}

	return web.ResponseOK(w, "news", article)
}

// GetLocationFeed is the main feed service. Will return a fixed number of articles. Locations are needed as query
// params in l1 and l2. If those are not provided, default ones are added.
func (h *Handler) GetLocationFeed(w http.ResponseWriter, r *http.Request) error {
	l1 := r.URL.Query().Get("l1")
	l2 := r.URL.Query().Get("l2")

	feed, err := h.Service.GetFeed(l1, l2)

	if err != nil {
		return err
	}

	return web.ResponseOK(w, "feed", feed)
}
