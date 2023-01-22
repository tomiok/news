package collector

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Storage interface {
	saveArticle(a Article) (*Article, error)
}

type storage struct {
	*sql.DB
}

func newStorage(url string) (*storage, error) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &storage{
		DB: db,
	}, nil
}

func (s *storage) saveArticle(a Article) (*Article, error) {
	res, err := s.Exec("insert into articles (title, description, content, link, country, location, lang, pub_date) values (?,?,?,?,?,?,?,?)")

	if err != nil {
		return nil, fmt.Errorf("cannot save article: %v", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("cannot get last inserted ID for articles: %v", err)
	}
	a.ID = id

	for _, catID := range a.Categories {
		if _, err := s.Exec("insert into article_categories (article_id, category_id) values (?,?)", a.ID, catID); err != nil {
			log.Warn().Err(err).Msg("cannot save article_categories")
		}
	}

	return &a, nil
}
