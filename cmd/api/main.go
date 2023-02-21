package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
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
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	r.Get("/news/{articleUID}", unwrap(deps.collectorHandler.GetNews))
	r.Get("/feeds", unwrap(deps.collectorHandler.GetLocationFeed))

	r.Get("/", unwrap(deps.collectorHandler.Home))
}

func collect(deps *dependencies) {
	ticker := time.NewTicker(10 * time.Minute)
	for _ = range ticker.C {
		deps.AggregateJob.Do()
	}
}
