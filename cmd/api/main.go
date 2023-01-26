package main

import (
	"fmt"
	"news/internal/collector"
	"time"
)

func main() {
	now := time.Now()
	collector.NewJob().Do()
	fmt.Println(time.Since(now))
}
