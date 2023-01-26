package main

import (
	"fmt"
	"news/internal/collector"
	"time"
)

const mysqlURI = "localhost"

func main() {
	now := time.Now()
	job, err := collector.NewJob("root:root@tcp(localhost:3306)/db")
	if err != nil {
		panic(err)
	}

	job.Do()
	fmt.Println(time.Since(now))
}
