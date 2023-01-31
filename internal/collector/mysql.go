package collector

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

const (
	maxOpenConnections = 10
	maxIdleConnections = 10
)

type Storage interface {
	saveArticle(a Article) (*Article, error)
}

type SQLStorage struct {
	*sql.DB
}

func NewStorage(url string) (*SQLStorage, error) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &SQLStorage{
		DB: db,
	}, nil
}

func (s *SQLStorage) saveArticle(a Article) (*Article, error) {
	res, err := s.Exec("insert into articles (title, uid, description, content, link, country, location, lang, pub_date, saved_at) values (?,?,?,?,?,?,?,?,?,?)",
		a.Title, a.UID, a.Description, a.Content, a.Link, a.Country, a.Location, a.Lang, a.PubDate, a.SavedAt)

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
