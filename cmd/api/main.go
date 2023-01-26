package main

import (
	"fmt"
	"news/internal/collector"
	"time"
)

func main() {
	now := time.Now()
	job, _ := collector.NewJob("")
	job.Do()
	fmt.Println(time.Since(now))
}
