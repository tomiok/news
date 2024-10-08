package main

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"news/platform/web"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	*http.Server
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (s *server) Start() {
	log.Info().Msgf("server is running on port %s", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("closed server error %s", err.Error())
		}
	}()
	s.gracefulShutdown()
}

func (s *server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT)
	sig := <-quit
	log.Info().Msgf("server is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		log.Error().Msgf("could not gracefully shutdown the server %s", err.Error())
	}
	log.Info().Msg("server stopped")
}

type webHandler func(w http.ResponseWriter, r *http.Request) error

func unwrap(f webHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			requestID := middleware.GetReqID(r.Context())
			log.Error().Caller(1).Err(err).Str("RequestID", requestID).Msg("cannot process request")
			web.ReturnErr(w, err)
			return
		}
	}
}
