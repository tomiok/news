// Package collector is the responsible to hit RSS feeds and save into database.
package collector

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const hours24 = 86400000 //24 hours in millis

const rssTimeout = 5 * time.Second

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
	Country     string // ISO code for the country AR, UY, BR...
	Location    string // Specific location for a specific site.
	PubDate     int64
	Categories  []string
}

// Article is the struct to save in the DB. The categories are curated, and we can save them safe.
type Article struct {
	ID          int64
	UID         string // for link generation proposes.
	Title       string
	Description string
	Content     string
	Country     string
	Location    string

	Lang string
	Link string

	PubDate int64
	SavedAt int64

	Categories []int // we have the category ids here.
}

// RSSCollector the RSS implementation of the Collector interface.
type RSSCollector struct {
	Parser *gofeed.Parser
}

// NewService is for web API only and returns *Service and an Error.
func NewService(url string, views []string) (*Service, error) {
	if views == nil || len(views) == 0 {
		return nil, errors.New("views are nil or empty")
	}

	storage := NewStorage(url)

	return &Service{
		Storage: storage,
	}, nil
}

// NewCollector returns a *Collector.
func NewCollector() *RSSCollector {
	p := gofeed.NewParser()
	p.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	return &RSSCollector{
		Parser: p,
	}
}

func (r *RSSCollector) Collect(ctx context.Context, site Site) ([]RawArticle, error) {
	now := time.Now().UnixMilli()

	_ctx, cancel := context.WithTimeout(ctx, rssTimeout)
	defer cancel()
	feed, err := r.Parser.ParseURLWithContext(site.URL, _ctx)

	if err != nil {
		return nil, fmt.Errorf("cannot parse feed for URL %s - %v", site.URL, err)
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

type Site struct {
	URL          string
	MainCategory string
	HasContent   bool
	Country      string
	Location     string
}

type SiteScanner struct {
	file string
}

func NewSiteScanner() *SiteScanner {
	return &SiteScanner{
		file: "internal/collector/sites.csv",
	}
}

func (s SiteScanner) Scan() []Site {
	f, err := os.Open(s.file)
	if err != nil {
		log.Fatal("unable to read input file "+s.file, err)
	}

	defer func() {
		_ = f.Close()
	}()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("unable to parse file as CSV for "+s.file, err)
	}

	result := make([]Site, 0, len(records))
	for _, lines := range records {
		_url := lines[0]
		mainCategory := lines[1]
		hasContent, _ := strconv.ParseBool(lines[2])
		country := lines[3]
		location := lines[4]
		result = append(result, Site{
			URL:          _url,
			MainCategory: mainCategory,
			HasContent:   hasContent,
			Country:      country,
			Location:     location,
		})
	}

	return result
}
