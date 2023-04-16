package main

import (
	"news/internal/feed"
	collectorHandler "news/internal/feed/handler"
	"os"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	mysqlURI  = "root:@tcp(localhost:3306)/news"
	localhost = "localhost:" + portLocal
)

type dependencies struct {
	AggregateJob     *feed.JobContainer
	collectorHandler *collectorHandler.Handler

	Port        string
	Environment string // which env is the program running.

	mySqlURI string
}

func newDeps() *dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)
	host := getVar("HOST", localhost)
	dbURI := getVar("DB_URI", mysqlURI)

	_job, err := feed.NewJob(host, dbURI)

	if err != nil {
		panic(err)
	}

	_collectorHandler, err := collectorHandler.New(dbURI)

	if err != nil {
		panic(err)
	}

	return &dependencies{
		AggregateJob:     _job,
		collectorHandler: _collectorHandler,

		Environment: env,
		Port:        port,
		mySqlURI:    dbURI,
	}
}

func getVar(key, defValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defValue
	}
	return val
}
