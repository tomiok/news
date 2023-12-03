package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"news/internal/feed"
	collectorHandler "news/internal/feed/handler"
	"os"
	"strconv"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	dbHostLocal     = "localhost"
	dbNameLocal     = "news"
	dbUserLocal     = "news"
	dbPasswordLocal = "news"
	dbPortLocal     = "5432"
	templateCache   = false
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

	dbName := getVar("POSTGRES_DB", dbNameLocal)
	dbUser := getVar("POSTGRES_USER", dbUserLocal)
	dbPassword := getVar("POSTGRES_PASSWORD", dbPasswordLocal)

	dbPort := getVar("DB_PORT", dbPortLocal)
	dbHost := getVar("DB_HOST", dbHostLocal)

	tempCache := getBoolVar("CACHE", templateCache)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	_storage := feed.NewStorage(dsn)
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
