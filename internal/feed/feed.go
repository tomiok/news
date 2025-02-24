package feed

import (
	"github.com/rs/zerolog/log"
	"time"
)

const (
	Argentina = "argentina"
	CABA      = "caba"
)

// Storage will interact with the DB.
type Storage interface {
	SaveArticle(a Article) (Article, error)
	GetArticleByUID(uid string) (Article, error)

	GetDBFeed(q string) ([]Article, error)

	GetSites() ([]Site, error)
}

// GetNewsByUID give a UID (stored in DB) return an *Article.
func (s *Service) GetNewsByUID(uid string) (Article, error) {
	article, err := s.Storage.GetArticleByUID(uid)

	if err != nil {
		return Article{}, err
	}

	return article, nil
}

// GetFeed will return a slice of articles. A pair of locations will be given, if is empty, a default one will be added.
func (s *Service) GetFeed(q string) ([]Article, error) {
	feed, err := s.Storage.GetDBFeed(q)

	if err != nil {
		return nil, err
	}

	return feed, nil
}

func Collect(d time.Duration, job JobAggregator) {
	ticker := time.NewTicker(d)
	for _ = range ticker.C {
		now := time.Now()
		job.Do()

		log.Info().Msgf("job duration: %s", time.Since(now))
	}
}
