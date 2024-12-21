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
	log.Info().Msgf("running job at: %s", time.Now().Format(time.RFC3339))
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
	close(done)
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
		go func(rawArticle RawArticle) {
			defer wg.Done()
			uid := a.GenerateID()
			_, err := a.storage.saveArticle(Article{
				Title:       rawArticle.Title,
				UID:         uid,
				Description: rawArticle.Description,
				Content:     template.HTML(rawArticle.Content),
				RawContent:  rawArticle.RawContent,
				Country:     rawArticle.Country,
				Location:    rawArticle.Location,
				PubDate:     rawArticle.PubDate,
				Link:        createLink(rawArticle.Title, uid),
				Source:      rawArticle.Source,
				SavedAt:     time.Now().UnixMilli(),
				Lang:        getLang(rawArticle.Country),
				Categories:  rawArticle.Categories,
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
		Argentina: langSpanish,
	}

	return m[country]
}

func createLink(title, uid string) string {
	return fmt.Sprintf("/news/%s/%s?permaLink=true", slug.MakeLang(title, "es"), uid)
}
