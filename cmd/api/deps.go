package main

import (
	"github.com/rs/zerolog/log"
	"news/internal/feed"
	collectorHandler "news/internal/feed/handler"
	"os"
	"strconv"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	mysqlURI      = "tomi:tomi@tcp(localhost:3306)/news_api_dev"
	templateCache = false
)

type dependencies struct {
	AggregateJob     *feed.JobContainer
	collectorHandler *collectorHandler.Handler

	Port        string
	Environment string // which env is the program running.

	CacheTemplate bool //template is going to be cached (only true in prod).
}

func newDeps() *dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)
	dbURI := getVar("DB_URI", mysqlURI)
	tempCache := getBoolVar("CACHE", templateCache)

	_storage := feed.NewStorage(dbURI)
	_job, err := feed.NewJob(_storage)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	_collectorHandler, err := collectorHandler.New(_storage, tempCache)

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

func getBoolVar(key string, defValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defValue
	}

	res, err := strconv.ParseBool(val)
	if err != nil {
		return defValue
	}
	return res
}

func getVar(key, defValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defValue
	}
	return val
}
