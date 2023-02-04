package main

import (
	"news/internal/collector"
	collectorHandler "news/internal/collector/handler"
	"os"
)

const (
	envLocal  = "local"
	portLocal = "9000"

	mysqlURI  = "root:root@tcp(localhost:3306)/db"
	localhost = "localhost:" + portLocal
)

type dependencies struct {
	AggregateJob     *collector.AggregateJob
	collectorHandler *collectorHandler.Handler

	Port        string
	Environment string // which env is the program running.

	mySqlURI string
	Host     string
}

func newDeps() *dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)
	host := getVar("HOST", localhost)
	dbURI := getVar("DB_URI", mysqlURI)

	_job, err := collector.NewJob(host, dbURI)

	if err != nil {
		panic(err)
	}

	_collectorHandler := collectorHandler.New(dbURI)

	return &dependencies{
		AggregateJob:     _job,
		collectorHandler: _collectorHandler,

		Environment: env,
		Port:        port,
		Host:        host,
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
