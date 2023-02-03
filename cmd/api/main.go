package main

import (
	"fmt"
	"net/http"
	"news/internal/collector"
	"time"
)

const mysqlURI = "root:root@tcp(localhost:3306)/db"

func main() {
	now := time.Now()
	job, err := collector.NewJob("localhost:8080", mysqlURI)
	if err != nil {
		panic(err)
	}

	job.Do()
	collector.Print()
	fmt.Println(time.Since(now))
}

func run() {
	srv := &http.Server{
		Addr: ":" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	server := server{srv}
	server.Start()
}
