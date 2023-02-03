package main

import (
	collectorHandler "news/internal/collector/handler"
	"os"
)

const (
	envLocal  = "local"
	portLocal = "9000"
)

type dependencies struct {
	collectorHandler collectorHandler.Handler

	Port        string
	Environment string // which env is the program running.
}

func newDeps() *dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)

	return &dependencies{
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
