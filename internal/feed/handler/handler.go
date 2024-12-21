package handler

import (
	"net/http"
	"news/internal/feed"
	"news/platform/web"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Handler will carry the services logic to the web layer.
type Handler struct {
	*feed.Service
	Cache bool
}

// New returns a *Handler if the service is created OK, otherwise an error.
func New(storage feed.Storage, cache bool) (*Handler, error) {
	service, err := feed.NewService(storage)

	if err != nil {
		return nil, err
	}

	return &Handler{
		Service: service,
		Cache:   cache,
	}, nil
}

// GetNews is a web handler that will return a news by the UID, provided in the path param.
// I.E. /news/some-title-here/i2W6uBcVvyDPiSCn4iMo8G
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

	return web.TemplateRender(w, "news.page.tmpl", &web.TemplateData{
		Article: &article,
	}, h.Cache)
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) error {
	l := r.URL.Query().Get("l1")
	var locations []string
	if l == "" {
		locations = []string{feed.Argentina, feed.CABA}
	} else {
		locations = strings.Split(l, ",")
	}

	articles, err := h.Service.GetFeed(locations...)
	if err != nil {
		return err
	}

	return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{
		Locations: strings.Join(locations, ","),
		Articles:  articles,
	}, h.Cache)
}

const maxLocations = 3

func (h *Handler) FeedsLookup(w http.ResponseWriter, r *http.Request) error {
	l := r.URL.Query().Get("l")
	locations := strings.Split(l, ",")

	if len(locations) > maxLocations {
		locations = locations[0:2]
	}

	articles, err := h.Service.GetFeed(locations...)
	if err != nil {
		return err
	}

	return web.TemplateRender(w, "feed.news.page.tmpl", &web.TemplateData{
		Locations: strings.Join(locations, ","),
		Articles:  articles,
	}, h.Cache)
}
