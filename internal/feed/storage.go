package feed

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	maxOpenConnections = 25
	maxIdleConnections = 20
)

// Storage will interact with the DB.
type Storage interface {
	saveArticle(a Article) (Article, error)
	getArticleByUID(uid string) (Article, error)

	GetDBFeed(locs ...string) ([]Article, error)

	GetSites() ([]Site, error)
}

type SQLStorage struct {
	*sql.DB
}

func NewStorage(url string) *SQLStorage {
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal().Err(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &SQLStorage{
		DB: db,
	}
}

func (s *SQLStorage) saveArticle(a Article) (Article, error) {
	res, err := s.Exec(`insert into articles 
    (title, uid, description, content, raw_content, link, country, location, lang, source, pub_date, saved_at,categories) 
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		a.Title,
		a.UID,
		a.Description,
		a.Content,
		a.RawContent,
		a.Link,
		strings.ToLower(a.Country),
		strings.ToLower(a.Location),
		a.Lang,
		a.Source,
		a.PubDate,
		a.SavedAt,
		pq.Array(a.Categories))

	if err != nil {
		return Article{}, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return Article{}, fmt.Errorf("cannot get last inserted ID for articles: %w", err)
	}
	a.ID = id

	return a, nil
}

func (s *SQLStorage) getArticleByUID(uid string) (Article, error) {
	var article Article
	row := s.QueryRow("select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.country, a.location, a.lang, a.source, a.pub_date, a.categories from articles a where a.uid=$1", uid)
	var categories []string
	err := row.Scan(
		&article.ID,
		&article.UID,
		&article.Title,
		&article.Description,
		&article.Content,
		&article.RawContent,
		&article.Country,
		&article.Location,
		&article.Lang,
		&article.Source,
		&article.PubDate,
		pq.Array(&categories),
	)

	if err != nil {
		return article, fmt.Errorf("cannot get articles %w", err)
	}

	article.Categories = categories
	return article, nil
}

const defSize = 50

func (s *SQLStorage) GetDBFeed(locations ...string) ([]Article, error) {
	back48Hours := time.Now().Add(-time.Hour * 48).UnixMilli()
	if locations == nil || len(locations) == 0 {
		return nil, errors.New("locations are nil or empty")
	}

	rows, err := s.Query("select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.link, a.country, a.location, a.lang, a.pub_date, a.categories from articles a where a.location in ($1) and a.pub_date >= $2 ORDER BY RANDOM() limit 50",
		strings.ToLower(locations[0]), back48Hours,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var categories []string
	result := make([]Article, 0, defSize)
	for rows.Next() {
		var article Article
		err = rows.Scan(
			&article.ID,
			&article.UID,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.RawContent,
			&article.Link,
			&article.Country,
			&article.Location,
			&article.Lang,
			&article.PubDate,
			pq.Array(&categories),
		)
		if err != nil {
			log.Error().Err(err).Msg("cannot read article")
			continue
		}
		article.SinceMinutes()
		article.Categories = categories
		result = append(result, article)
	}

	return result, nil
}

func (s *SQLStorage) GetSites() ([]Site, error) {
	rows, err := s.Query("select url, category, has_content, country, location from sites")

	if err != nil {
		return nil, err
	}

	var result []Site
	for rows.Next() {
		var site Site
		err = rows.Scan(&site.URL, &site.MainCategory, &site.HasContent, &site.Country, &site.Location)
		if err != nil {
			log.Error().Err(err)
		}
		result = append(result, site)
	}

	return result, nil
}
