package collector

import (
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Job will be the abstraction of get the read the CSV, get the RSS, sanitize and save in Database.
type Job interface {
	Do()
}

type AggregateJob struct {
	Scanner   Scanner
	Collector Collector
	Sanitizer Sanitizer
	Storage   Storage
}

func NewJob(mysqlURI string) (*AggregateJob, error) {
	storage, err := NewStorage(mysqlURI)
	if err != nil {
		return nil, err
	}

	return &AggregateJob{
		Collector: NewCollector(),
		Scanner:   NewSiteScanner(),
		Sanitizer: NewSanitizer(),
		Storage:   storage,
	}, nil
}

func (a *AggregateJob) Do() {
	// declare channels
	chSites := make(chan Site)
	chArticles := make(chan RawArticle)
	transformedCh := make(chan RawArticle)
	done := make(chan struct{})
	go a.getSites(chSites)
	go a.getRawArticles(chSites, chArticles)
	go a.Sanitize(chArticles, transformedCh)
	go a.Save(transformedCh, done)
	<-done
}

func (a *AggregateJob) GenerateID() string {
	return shortuuid.New()
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

func (a *AggregateJob) Save(ch chan RawArticle, done chan struct{}) {
	var wg sync.WaitGroup
	for article := range ch {
		wg.Add(1)
		go func(article RawArticle) {
			defer wg.Done()
			_, err := a.Storage.saveArticle(Article{
				Title:       article.Title,
				UID:         a.GenerateID(),
				Description: article.Description,
				Content:     article.Content,
				Country:     article.Country,
				Location:    article.Location,
				PubDate:     article.PubDate,
				SavedAt:     time.Now().UnixMilli(),
				Lang:        getLang(article.Country),
				Categories:  []int{},
			})
			if err != nil {
				log.Warn().Err(err).Msg("")
			}
			saveCategory(article)
		}(article)
	}
	wg.Wait()
	done <- struct{}{}
}

func getLang(country string) string {
	var m = map[string]string{
		countryAR: langSpanish,
	}

	return m[country]
}

var _m = sync.Map{}

func saveCategory(article RawArticle) {
	for _, category := range article.Categories {
		_m.Store(category, struct{}{})
	}
}

func Print() {
	_m.Range(func(key, value any) bool {
		fmt.Println(key)
		return true
	})
}
