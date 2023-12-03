package feed

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
	"html/template"
	"sync"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/rs/zerolog/log"
)

// Lock is a mechanism to run only once the Job that get the feeds and save it in the database.
type Lock struct {
	IsLocked  bool
	Timestamp int64
}

// JobAggregator will be the abstraction of get the read the CSV, get the RSS, sanitize and save in Database.
type JobAggregator interface {
	Do()
}

type JobContainer struct {
	scanner   Scanner
	collector Collector
	sanitizer Sanitizer
	storage   Storage
}

func NewJob(storage Storage) (*JobContainer, error) {
	return &JobContainer{
		collector: newCollector(),
		scanner:   newSiteScanner(storage),
		sanitizer: newSanitizer(),
		storage:   storage,
	}, nil
}

func (a *JobContainer) Do() {
	lock, err := a.storage.AcquireLock()
	if err != nil {
		log.Warn().Err(err).Msg("collector is locked")
		return
	}

	if time.Now().Sub(time.UnixMilli(lock.Timestamp)).Abs() < time.Hour {
		log.Warn().Msg("last run was less than 1 hour")
		return
	}

	log.Info().Msgf("running job: %s", time.UnixMilli(lock.Timestamp))
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

func (a *JobContainer) getRawArticles(sitesCh chan Site, out chan RawArticle) {
	var wg sync.WaitGroup
	for site := range sitesCh {
		wg.Add(1)
		go func(site Site) {
			defer wg.Done()
			articles, err := a.collector.Collect(context.Background(), site)
			if err != nil {
				log.Warn().Err(err).Msgf("cannot found articles for %s", site.URL)
			} else {
				for _, article := range articles {
					out <- article
				}
			}
		}(site)
	}
	wg.Wait()
	close(out)
}

func (a *JobContainer) Sanitize(articlesCh, out chan RawArticle) {
	var wg sync.WaitGroup
	for rawArt := range articlesCh {
		wg.Add(1)
		title, desc, content, rawContent := a.sanitizer.Apply(rawArt.Title, rawArt.Description, rawArt.Content)
		go func(t, d, c, rc string, rawArt RawArticle) {
			wg.Done()
			out <- RawArticle{
				Title:       t,
				Description: d,
				Content:     c,
				RawContent:  rc,
				Country:     rawArt.Country,
				Location:    rawArt.Location,
				PubDate:     rawArt.PubDate,
				Categories:  rawArt.Categories,
			}
		}(title, desc, content, rawContent, rawArt)
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
				Content:     template.HTML(article.Content),
				RawContent:  article.RawContent,
				Country:     article.Country,
				Location:    article.Location,
				PubDate:     article.PubDate,
				Link:        createLink(article.Title, uid),
				Source:      article.Source,
				SavedAt:     time.Now().UnixMilli(),
				Lang:        getLang(article.Country),
				Categories:  []int{},
			})
			if err != nil {
				_, ok := err.(*mysql.MySQLError)
				if !ok {
					return
				}
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

func createLink(title, uid string) string {
	return fmt.Sprintf("/news/%s/%s?permaLink=true", slug.MakeLang(title, "es"), uid)
}
