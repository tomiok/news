package collector

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

var arr = []string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "pre", "code", "blockquote", "article", "section"}

type Job interface {
	Do()
}

type AggregateJob struct {
	Scanner   Scanner
	Collector Collector
	Storage   Storage
}

func NewJob() *AggregateJob {
	return &AggregateJob{
		Collector: NewCollector(),
		Scanner:   NewSiteScanner(),
	}
}

func (a *AggregateJob) Do() {
	sites := a.Scanner.Scan()
	chSites := make(chan Site)
	chArticles := make(chan RawArticle)

	go func() {
		for _, site := range sites {
			chSites <- site
		}
		close(chSites)
	}()

	go a.do(chSites, chArticles)
	Print(chArticles)
}

func (a *AggregateJob) do(sitesCh chan Site, articlesCh chan RawArticle) {
	for site := range sitesCh {
		articles, err := a.Collector.Collect(context.Background(), site)
		if err != nil {
			log.Warn().Err(err).Msgf("cannot found articles for %s", site.URL)
		}
		for _, article := range articles {
			articlesCh <- article
		}
		//close(articlesCh)

	}
}

func Print(articlesCh chan RawArticle) {
	for a := range articlesCh {
		fmt.Println(a.String())
	}
}
