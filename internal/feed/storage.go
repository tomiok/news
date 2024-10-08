package feed

import (
	"database/sql"
	"errors"
	"fmt"
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
	saveArticle(a Article) (*Article, error)
	getArticleByUID(uid string) (*Article, error)

	GetDBFeed(locs ...string) ([]Article, error)

	AcquireLock() (*Lock, error)

	GetSites() ([]Site, error)
}

type SQLStorage struct {
	*sql.DB
}

func NewStorage(url string) *SQLStorage {
	fmt.Println(url)
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

func (s *SQLStorage) saveArticle(a Article) (*Article, error) {
	res, err := s.Exec("insert into articles (title, uid, description, content, raw_content, link, country, location, lang, source, pub_date, saved_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
		a.Title, a.UID, a.Description, a.Content, a.RawContent, a.Link, a.Country, a.Location, a.Lang, a.Source, a.PubDate, a.SavedAt)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("cannot get last inserted ID for articles: %w", err)
	}
	a.ID = id

	for _, catID := range a.Categories {
		if _, err := s.Exec("insert into article_categories (article_id, category_id) values ($1,$2)", a.ID, catID); err != nil {
			log.Warn().Err(err).Msg("cannot save article_categories")
		}
	}
	return &a, nil
}

func (s *SQLStorage) getArticleByUID(uid string) (*Article, error) {
	var article Article
	row := s.QueryRow("select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.country, a.location, a.lang, a.source, a.pub_date from articles a where a.uid=$1", uid)
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
	)

	if err != nil {
		return nil, fmt.Errorf("cannot get articles %w", err)
	}

	return &article, nil
}

const defSize = 50

func (s *SQLStorage) GetDBFeed(locations ...string) ([]Article, error) {
	oneDay := time.Now().Add(-time.Hour * 48).UnixMilli()
	if locations == nil || len(locations) == 0 {
		return nil, errors.New("locations are nil or empty")
	}
	rows, err := s.Query("select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.link, a.country, a.location, a.lang, a.pub_date from articles a where a.location in ($1,$2) and a.pub_date >= $3 ORDER BY RANDOM() limit 50",
		locations[0], locations[1], oneDay,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

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
		)
		if err != nil {
			log.Error().Err(err).Msg("cannot read article")
			continue
		}
		article.SinceMinutes()
		result = append(result, article)
	}

	return result, nil
}

func (s *SQLStorage) AcquireLock() (*Lock, error) {
	tx, err := s.Begin()
	if err != nil {
		return nil, err
	}
	var countID int
	err = tx.QueryRow("select id from feed_lock order by id asc limit 1").Scan(&countID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = tx.Exec("insert into feed_lock (is_locked, timestamp) values($1,$2)", true, time.Now().UnixMilli())
			return &Lock{
				IsLocked:  true,
				Timestamp: time.Now().Add(-2 * time.Hour).UnixMilli(),
			}, tx.Commit()
		}
		return nil, err
	}
	var ts int64
	var isLocked bool
	err = tx.QueryRow("select is_locked, timestamp from feed_lock order by id asc limit 1").Scan(&isLocked, &ts)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if time.Now().Sub(time.UnixMilli(ts)).Abs() < time.Hour {
		return nil, errors.New("already locked")
	}

	_, err = tx.Exec("insert into feed_lock (is_locked, timestamp) values($1,$2)", true, time.Now().UnixMilli())
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &Lock{IsLocked: false, Timestamp: ts}, tx.Commit()
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
