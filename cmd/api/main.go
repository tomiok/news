package main

import "news/internal/collector"

func main() {
	collector.NewJob().Do()
}
