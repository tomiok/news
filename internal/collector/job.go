package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"news/internal/collector"
	"strings"
	"sync"
)

var arr = []string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "pre", "code", "blockquote", "article", "section"}

type Job interface {
	Do()
}

type AggregateJob struct {
	Scanner   collector.Scanner
	Collector collector.Collector
	Storage   collector.Storage
}

func NewJob() *AggregateJob {
	return &AggregateJob{
		Collector: collector.NewCollector(),
		Scanner:   collector.NewSiteScanner(),
	}
}

func (a *AggregateJob) Do() {
	// declare channels
	chSites := make(chan collector.Site)
	chArticles := make(chan collector.RawArticle)

	go a.getSites(chSites)
	go a.getRawArticles(chSites, chArticles)

	Print(chArticles)
}

func (a *AggregateJob) getSites(chSites chan collector.Site) {
	sites := a.Scanner.Scan()
	for _, site := range sites {
		chSites <- site
	}
	close(chSites)
}

func (a *AggregateJob) getRawArticles(sitesCh chan collector.Site, articlesCh chan collector.RawArticle) {
	var wg sync.WaitGroup
	for site := range sitesCh {
		wg.Add(1)
		go func(site collector.Site) {
			articles, err := a.Collector.Collect(context.Background(), site)
			if err != nil {
				log.Warn().Err(err).Msgf("cannot found articles for %s", site.URL)
			}
			for _, article := range articles {
				articlesCh <- article
			}
			wg.Done()
		}(site)
		wg.Wait()
	}
	close(articlesCh)
}

func Print(articlesCh chan collector.RawArticle) {
	for a := range articlesCh {
		fmt.Println(strings.TrimSpace(a.Title))
	}
}
