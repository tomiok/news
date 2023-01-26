package collector

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"sync"
)

type Job interface {
	Do()
}

type AggregateJob struct {
	Scanner   Scanner
	Collector Collector
	Sanitizer Sanitizer
	Storage   Storage
}

func NewJob() *AggregateJob {
	return &AggregateJob{
		Collector: NewCollector(),
		Scanner:   NewSiteScanner(),
		Sanitizer: NewSanitizer(),
	}
}

func (a *AggregateJob) Do() {
	// declare channels
	chSites := make(chan Site)
	chArticles := make(chan RawArticle)
	transformedCh := make(chan RawArticle)

	go a.getSites(chSites)
	go a.getRawArticles(chSites, chArticles)
	go a.Sanitize(chArticles, transformedCh)
	Print(transformedCh)
}

func (a *AggregateJob) getSites(chSites chan Site) {
	sites := a.Scanner.Scan()
	for _, site := range sites {
		chSites <- site
	}
	close(chSites)
}

func (a *AggregateJob) getRawArticles(sitesCh chan Site, articlesCh chan RawArticle) {
	var wg sync.WaitGroup
	for site := range sitesCh {
		wg.Add(1)
		go func(site Site) {
			articles, err := a.Collector.Collect(context.Background(), site)
			if err != nil {
				log.Warn().Err(err).Msgf("cannot found articles for %s", site.URL)
			}
			for _, article := range articles {
				articlesCh <- article
			}
			wg.Done()
		}(site)
	}
	wg.Wait()
	close(articlesCh)
}

func (a *AggregateJob) Sanitize(articlesCh, out chan RawArticle) {
	var wg sync.WaitGroup
	for rawArt := range articlesCh {
		wg.Add(1)
		title, desc, content := a.Sanitizer.Apply(rawArt.Title, rawArt.Description, rawArt.Content)
		go func(t, d, c string, rawArt RawArticle) {
			art := RawArticle{
				Title:       t,
				Description: d,
				Content:     c,
				Country:     rawArt.Country,
				Location:    rawArt.Location,
				PubDate:     rawArt.PubDate,
				Categories:  rawArt.Categories,
			}
			out <- art
			wg.Done()
		}(title, desc, content, rawArt)
	}
	wg.Wait()
	close(out)
}

func Print(articlesCh chan RawArticle) {
	for a := range articlesCh {
		fmt.Println(a.Title)
	}
}
