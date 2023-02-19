package feed

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/rs/zerolog/log"
)

// JobAggregator will be the abstraction of get the read the CSV, get the RSS, sanitize and save in Database.
type JobAggregator interface {
	Do()
}

type JobContainer struct {
	scanner   Scanner
	collector Collector
	sanitizer Sanitizer
	storage   Storage

	Host string
}

func NewJob(host, mysqlURI string) (*JobContainer, error) {
	storage := NewStorage(mysqlURI)

	return &JobContainer{
		collector: newCollector(),
		scanner:   newSiteScanner(),
		sanitizer: newSanitizer(),
		storage:   storage,
		Host:      host,
	}, nil
}

func (a *JobContainer) Do() {
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

func (a *JobContainer) GenerateID() string {
	return shortuuid.New()
}

func (a *JobContainer) getSites(chSites chan Site) {
	sites := a.scanner.Scan()
	for _, site := range sites {
		chSites <- site
	}
	close(chSites)
}

func (a *JobContainer) getRawArticles(sitesCh chan Site, articlesCh chan RawArticle) {
	var wg sync.WaitGroup
	for site := range sitesCh {
		wg.Add(1)
		go func(site Site) {
			defer wg.Done()
			articles, err := a.collector.Collect(context.Background(), site)
			if err != nil {
				log.Warn().Err(err).Msgf("cannot found articles for %s", site.URL)
				return
			}
			for _, article := range articles {
				articlesCh <- article
			}
		}(site)
	}
	wg.Wait()
	close(articlesCh)
}

func (a *JobContainer) Sanitize(articlesCh, out chan RawArticle) {
	var wg sync.WaitGroup
	for rawArt := range articlesCh {
		wg.Add(1)

		title, desc, content := a.sanitizer.Apply(rawArt.Title, rawArt.Description, rawArt.Content)
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

func (a *JobContainer) Save(ch chan RawArticle, done chan struct{}) {
	var wg sync.WaitGroup
	for article := range ch {
		wg.Add(1)
		go func(article RawArticle) {
			defer wg.Done()
			uid := a.GenerateID()
			_, err := a.storage.saveArticle(Article{
				Title:       article.Title,
				UID:         uid,
				Description: article.Description,
				Content:     article.Content,
				Country:     article.Country,
				Location:    article.Location,
				PubDate:     article.PubDate,
				Link:        createLink(uid),
				SavedAt:     time.Now().UnixMilli(),
				Lang:        getLang(article.Country),
				Categories:  []int{},
			})
			if err != nil {
				log.Warn().Err(err).Msg("")
			}
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

func createLink(uid string) string {
	link := fmt.Sprintf("/news/%s?permaLink=true", uid)
	return link
}
