package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"news/internal/feed"
	"news/platform/web"
	"strconv"
)

// Handler will carry the services logic to the web layer.
type Handler struct {
	*feed.Service
	Cache bool

	SecureCookie bool
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

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.Service.GetFeed("")
	if err != nil {
		return err
	}

	logged := r.URL.Query().Get("loggedIn")
	justLogged, _ := strconv.ParseBool(logged)

	c, err := r.Cookie("auth_token")
	if err != nil {
		log.Error().Msg("cannot read cookie")
		return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{
			Articles:     articles,
			JustLoggedIn: justLogged,
		}, h.Cache)
	}

	var token string
	var isLogged bool
	if len(c.Value) > 40 {
		token = c.Value
		isLogged = true
	}

	return web.TemplateRender(w, "home.page.tmpl", &web.TemplateData{
		Articles:     articles,
		JustLoggedIn: justLogged,
		Token:        token,
		IsLogged:     isLogged,
	}, h.Cache)
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

func (h *Handler) FeedsSearch(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")
	if q == "" {
		return web.TemplateRender(w, "feed.news.page.tmpl", &web.TemplateData{
			Articles: []any{},
		}, h.Cache)
	}

	articles, err := h.Service.GetFeed(q)
	if err != nil {
		return err
	}

	return web.TemplateRender(w, "feed.news.page.tmpl", &web.TemplateData{
		Articles: articles,
	}, h.Cache)
}
