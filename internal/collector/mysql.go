package collector

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

const (
	maxOpenConnections = 10
	maxIdleConnections
)

// Storage will interact with the DB.
type Storage interface {
	saveArticle(a Article) (*Article, error)
	getArticleByUID(uid string) (*Article, error)

	GetDBFeed(locs ...string) ([]Article, error)
}

type SQLStorage struct {
	*sql.DB
}

func NewStorage(url string) *SQLStorage {
	db, err := sql.Open("mysql", url)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &SQLStorage{
		DB: db,
	}
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

func (s *SQLStorage) getArticleByUID(uid string) (*Article, error) {
	var article Article
	row := s.QueryRow("select a.id, a.uid, a.title, a.description, a.content, a.country, a.location, a.lang, a.pub_date from articles a where a.uid=?", uid)
	err := row.Scan(
		&article.ID,
		&article.UID,
		&article.Title,
		&article.Description,
		&article.Content,
		&article.Country,
		&article.Location,
		&article.Lang,
		&article.PubDate,
	)

	if err != nil {
		return nil, fmt.Errorf("cannot get articles %w", err)
	}

	return &article, nil
}

const defSize = 50

func (s *SQLStorage) GetDBFeed(locations ...string) ([]Article, error) {
	oneDay := time.Now().Add(-time.Hour * 24).UnixMilli()
	if locations == nil || len(locations) == 0 {
		return nil, errors.New("locations are nil or empty")
	}
	rows, err := s.Query("select a.id, a.uid, a.title, a.description, a.content, a.link, a.country, a.location, a.lang, a.pub_date from articles a where a.location in (?,?) and a.pub_date >= ? ORDER BY RAND() limit 50",
		locations[0], locations[1], oneDay,
	)

	if err != nil {
		return nil, err
	}

	result := make([]Article, 0, defSize)
	for rows.Next() {
		var article Article
		err := rows.Scan(
			&article.ID,
			&article.UID,
			&article.Title,
			&article.Description,
			&article.Content,
			&article.Link,
			&article.Country,
			&article.Location,
			&article.Lang,
			&article.PubDate,
		)
		if err != nil {
			log.Error().Err(err).Msg("cannot read article")
			continue
		}
		result = append(result, article)
	}

	return result, nil
}
