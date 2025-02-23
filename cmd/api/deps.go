package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"news/internal/feed"
	"news/internal/feed/handler"
	"news/internal/feed/storage"
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

type Dependencies struct {
	AggregateJob     *feed.JobContainer
	CollectorHandler *handler.Handler

	Port        string
	Environment string // which env is the program running.

	CacheTemplate bool //template is going to be cached (only true in prod).

	MigrationsDSN string
}

func NewDeps() *Dependencies {
	env := getVar("ENV", envLocal)
	port := getVar("PORT", portLocal)

	dbName := getVar("POSTGRES_DB", dbNameLocal)
	dbUser := getVar("POSTGRES_USER", dbUserLocal)
	dbPassword := getVar("POSTGRES_PASSWORD", dbPasswordLocal)

	dbPort := getVar("DB_PORT", dbPortLocal)
	dbHost := getVar("DB_HOST", dbHostLocal)

	tempCache := getBoolVar("CACHE", templateCache)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	migrationsDSN := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	_storage := storage.NewStorage(dsn)
	_job, err := feed.NewJob(_storage)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	_collectorHandler, err := handler.New(_storage, tempCache)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &Dependencies{
		AggregateJob:     _job,
		CollectorHandler: _collectorHandler,

		Environment:   env,
		Port:          port,
		MigrationsDSN: migrationsDSN,
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
