package feed

import "github.com/rs/zerolog/log"

// Scanner interface could fetch the data from some file, containing
// url, main-category, has-content,
type Scanner interface {
	Scan() []Site
}

// Site is expressed as a website to be scanned. Among the URL, some others values are there to help the collector
// when grabbing the data.
type Site struct {
	ID           int64
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
