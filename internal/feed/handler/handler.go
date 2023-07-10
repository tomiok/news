package handler

import (
	"net/http"
	"news/internal/feed"
	"news/platform/web"

	"github.com/go-chi/chi/v5"
)

// Handler will carry the services logic to the web layer.
type Handler struct {
	*feed.Service
}

// New returns a *Handler if the service is created OK, otherwise an error.
func New(storage feed.Storage) (*Handler, error) {
	service, err := feed.NewService(storage)

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

	_feed, locations, err := h.Service.GetFeed(l1, l2)

	if err != nil {
		return err
	}

	l := len(locations)
	for i := 0; i <= 2-l; i++ {
		locations = append(locations, "CABA")
	}

	return web.TemplateRender(w, "news.page.tmpl", &web.TemplateData{
		Articles:       _feed,
		FirstLocation:  locations[0],
		SecondLocation: locations[1],
	})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) error {
	l1 := r.URL.Query().Get("l1")
	l2 := r.URL.Query().Get("l2")

	if l1 == "" {
		l1 = feed.Argentina
	}

	if l2 == "" {
		l2 = feed.CABA
	}

	articles, locations, err := h.Service.GetFeed(l1, l2)
	if err != nil {
		return err
	}

	return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{
		FirstLocation:  locations[0],
		SecondLocation: locations[1],
		Articles:       articles,
	})
}
