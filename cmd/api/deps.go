package main

import (
	"news/internal/feed"
	collectorHandler "news/internal/feed/handler"
	"os"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	mysqlURI = "root:@tcp(localhost:3306)/news_dev"
)

//var connectionDB = fmt.Sprintf("%s:%s@tcp(%s:3306)/news_api_dev", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"))

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
	dbURI := getVar("DB_URI", mysqlURI)

	_job, err := feed.NewJob(feed.NewStorage(dbURI))

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
