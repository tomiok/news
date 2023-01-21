// Package collector is the responsible to hit RSS feeds and save into database.
package collector

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"os"
	"strconv"
	"time"
)

// Collector will be in charge of hit the RSS URL, and then, convert to an Article.
type Collector interface {
	// Collect simply return a list of Articles (items in RSS) parsed form a Site. The data is transformed in order
	// to save it easily in the Database.
	Collect(s Site) ([]Article, error)
}

// Article is the same as we can get in RSS feed.
type Article struct {
	ID          string
	Title       string
	Description string
	Country     string // ISO code for the country AR, UY, BR...
	Location    string // Specific location for a specific site.
	PubDate     int64
	Categories  []string
}

// RSSCollector the RSS implementation of the Collector interface.
type RSSCollector struct {
	Parser *gofeed.Parser
}

// NewCollector returns a *Collector.
func NewCollector() *RSSCollector {
	p := gofeed.NewParser()
	p.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	return &RSSCollector{
		Parser: p,
	}
}

func (r *RSSCollector) Collect(ctx context.Context, site Site) ([]Article, error) {
	_ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	feed, err := r.Parser.ParseURLWithContext(site.URL, _ctx)

	if err != nil {
		return nil, fmt.Errorf("cannot parse feed for URL %s - %v", site.URL, err)
	}

	result := make([]Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		var article Article
		article.Title = item.Title
		article.Description = item.Description
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

		result = append(result, article)
	}

	return result, nil
}

func (a Article) String() string {
	return fmt.Sprintf("Title: %s, desc: %s, cat: %s", a.Title, a.Description, a.Categories)
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

func (s SiteScanner) GetSites() []Site {
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
		log.Fatal("Unable to parse file as CSV for "+s.file, err)
	}

	result := make([]Site, 0, len(records))
	for _, lines := range records {
		b, _ := strconv.ParseBool(lines[2])
		result = append(result, Site{
			URL:          lines[0],
			MainCategory: lines[1],
			HasContent:   b,
			Country:      lines[3],
			Location:     lines[4],
		})
	}

	return result
}
