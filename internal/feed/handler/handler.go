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
		Article: article,
	}, h.Cache)
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

	articles, err := h.Service.GetFeed(l1, l2)
	if err != nil {
		return err
	}

	return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{
		FirstLocation:  l1,
		SecondLocation: l2,
		Articles:       articles,
	}, h.Cache)
}
