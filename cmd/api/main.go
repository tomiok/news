package main

import (
	"fmt"
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
