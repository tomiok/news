package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"news/internal/collector"
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

	now := time.Now()
	_, err := collector.NewJob("", mysqlURI)
	if err != nil {
		panic(err)
	}

	//job.Do()
	//collector.Print()
	fmt.Println(time.Since(now))

	routes(r, deps)
	serv := server{Server: srv}
	serv.Start()
}

func routes(r *chi.Mux, deps *dependencies) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	r.Get("/news/{articleUID}", unwrap(deps.collectorHandler.GetNews))
	r.Get("/feeds", unwrap(deps.collectorHandler.GetLocationFeed))
}
