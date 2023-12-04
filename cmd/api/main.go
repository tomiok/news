package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	run()
}

func run() {
	deps := newDeps()

	r := chi.NewRouter()
	srv := &http.Server{
		Addr: ":" + deps.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go collect(deps)

	routes(r, deps)
	serv := server{Server: srv}
	serv.Start()
}

func routes(r *chi.Mux, deps *dependencies) {
	r.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer, Cors(), middleware.Heartbeat("/ping"))

	r.Get("/news/{slug}/{articleUID}", unwrap(deps.collectorHandler.GetNews))

	r.Get("/", unwrap(deps.collectorHandler.Home))

	fileServer(r)
}

func collect(deps *dependencies) {
	ticker := time.NewTicker(1 * time.Hour)
	for _ = range ticker.C {
		now := time.Now()
		deps.AggregateJob.Do()

		log.Info().Msgf("job duration: %s", time.Since(now))
	}
}

func fileServer(r chi.Router) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fs(r, "/static", filesDir)
}

// fs conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("file server does not permit any URL parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		h := http.StripPrefix(pathPrefix, http.FileServer(root))
		h.ServeHTTP(w, r)
	})
}

func Cors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://web6am.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Content-Type-Options"},
		AllowCredentials: false,
		MaxAge:           500,
	})
}
