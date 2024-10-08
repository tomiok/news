// Package feed is the responsible to hit RSS feeds and save into database.
package feed

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// hours24, 24 hours in millis.
const hours24 = 86400000

const (
	countryAR   = "Argentina"
	langSpanish = "es_AR"
)

// Collector will be in charge of hit the RSS URL, and then, convert to an RawArticle.
type Collector interface {
	// Collect simply return a list of Articles (items in RSS) parsed form a Site. The data is transformed in order
	// to save it easily in the Database.
	Collect(ctx context.Context, s Site) ([]RawArticle, error)
}

// Service is a middleware for web API.
type Service struct {
	Storage
	views []string
}

// RawArticle is the same as we can get in RSS feed.
type RawArticle struct {
	HasContent  bool
	Title       string
	Description string // is like subtitle
	Content     string // content is the news itself. Some sites may don't have it.
	RawContent  string // do not use tags here.
	Country     string // ISO code for the country AR, UY, BR...
	Location    string // Specific location for a specific site.
	Source      string // add the actual web portal.
	PubDate     int64
	Categories  []string
}

// Article is the struct to save in the DB. The categories are curated, and we can save them safe.
type Article struct {
	ID          int64         `json:"id"`
	UID         string        `json:"-"` // for link generation proposes.
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Content     template.HTML `json:"content"`
	Country     string        `json:"country"`
	Location    string        `json:"location"`
	Source      string        `json:"-"`

	Lang string `json:"lang"`
	Link string `json:"link,omitempty"`

	PubDate int64  `json:"pub_date"`
	SavedAt int64  `json:"saved_at,omitempty"`
	Since   string `json:"since"`

	RawContent string `json:"parsed_content"`

	Categories []int `json:"categories,omitempty"` // we have the category ids here.
}

func (a *Article) SinceMinutes() {
	minutes := int(time.Since(time.UnixMilli(a.PubDate)).Minutes())

	if minutes <= 60 {
		a.Since = fmt.Sprintf("hace %d minutos", minutes)
	} else {
		hours := minutes / 60
		if hours == 1 {
			a.Since = "hace 1 hora"
		} else {
			a.Since = fmt.Sprintf("hace %d horas", hours)
		}
	}
}

// rssCollector the RSS implementation of the Collector interface.
type rssCollector struct {
	Parser *gofeed.Parser
}

// NewService is for web API only and returns *Service and an Error.
func NewService(storage Storage) (*Service, error) {
	return &Service{
		Storage: storage,
	}, nil
}

// newCollector returns a *Collector.
func newCollector() *rssCollector {
	p := gofeed.NewParser()
	p.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	return &rssCollector{
		Parser: p,
	}
}

func (r *rssCollector) Collect(ctx context.Context, site Site) ([]RawArticle, error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error().Msgf("recovered", rvr)
		}
	}()
	now := time.Now().UnixMilli()

	feed, err := r.Parser.ParseURLWithContext(site.URL, ctx)

	if err != nil {
		return nil, fmt.Errorf("cannot parse feed for URL %s - %w", site.URL, err)
	}

	result := make([]RawArticle, 0, len(feed.Items))
	for _, item := range feed.Items {
		var article RawArticle
		article.Title = item.Title
		if site.HasContent {
			article.Description = item.Description
			article.Content = item.Content
		} else {
			article.Content = item.Description
		}

		article.Country = site.Country
		article.Location = site.Location
		article.Source = site.URL
		if len(item.Categories) > 0 { // could be with 1 but still empty
			if item.Categories[0] != "" {
				article.Categories = item.Categories
			} else {
				article.Categories = []string{site.MainCategory}
			}
		} else {
			article.Categories = []string{site.MainCategory}
		}

		// means that is only a title for RSS. We don't want to save empty content.
		if strings.TrimSpace(article.Content) == "" {
			continue
		}

		var published int64

		if publishedAt := item.PublishedParsed; publishedAt != nil {
			published = publishedAt.UnixMilli()
		} else {
			published = time.Now().UnixMilli()
		}

		// do not save when are older than 24 hours
		if (now - published) > hours24 {
			continue
		}

		article.PubDate = published
		result = append(result, article)
	}

	return result, nil
}

func (a RawArticle) String() string {
	return fmt.Sprintf("Title: %s, desc: %s, content: %s, cat: %s", a.Title, a.Description, a.Content, a.Description)
}

// Scanner interface could fetch the data from some file, containing
// url, main-category, has-content,
type Scanner interface {
	Scan() []Site
}

// Site is expressed as a website to be scanned. Among the URL, some others values are there to help the collector
// when grabbing the data.
type Site struct {
	URL          string // the base URL of the RSS.
	MainCategory string // The main category added if the feed do not provide any other.
	HasContent   bool   // some RSS do not provide the content. Let's use the Description then.
	Country      string // Country of the site.
	Location     string // City or other location (province, state, etc)
}

type siteScanner struct {
	Storage
}

func newSiteScanner(storage Storage) *siteScanner {
	return &siteScanner{
		Storage: storage,
	}
}

func (s *siteScanner) Scan() []Site {
	sites, err := s.Storage.GetSites()

	if err != nil {
		log.Error().Err(err)
	}
	return sites
}
