package main

import (
	"github.com/rs/zerolog/log"
	"news/internal/feed"
	collectorHandler "news/internal/feed/handler"
	"os"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	mysqlURI = "tomi:tomi@tcp(localhost:3306)/news_api_dev"
)

type dependencies struct {
	AggregateJob     *feed.JobContainer
	collectorHandler *collectorHandler.Handler

	Port        string
	Environment string // which env is the program running.
}

func newDeps() *dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)
	dbURI := getVar("DB_URI", mysqlURI)

	_storage := feed.NewStorage(dbURI)
	_job, err := feed.NewJob(_storage)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	_collectorHandler, err := collectorHandler.New(_storage)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &dependencies{
		AggregateJob:     _job,
		collectorHandler: _collectorHandler,

		Environment: env,
		Port:        port,
	}
}

func getVar(key, defValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defValue
	}
	return val
}
