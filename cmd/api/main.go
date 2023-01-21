package main

import (
	"context"
	"fmt"
	"news/internal/collector"
)

func main() {
	scanner := collector.NewSiteScanner()

	values := scanner.GetSites()

	rssCollector := collector.NewCollector()

	res, _ := rssCollector.Collect(context.Background(), values[0])

	for _, v := range res {
		fmt.Println(v.String())
	}
}
