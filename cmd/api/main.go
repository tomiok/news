package main

import (
	"fmt"
	"news/internal/collector"
	"time"
)

const mysqlURI = "root:root@tcp(localhost:3306)/db"

func main() {
	now := time.Now()
	job, err := collector.NewJob(mysqlURI)
	if err != nil {
		panic(err)
	}

	job.Do()
	collector.Print()
	fmt.Println(time.Since(now))
}
