package feed

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

const (
	maxOpenConnections = 100
	maxIdleConnections
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
	db, err := sql.Open("mysql", url)

	if err != nil {
		log.Fatal().Err(err)
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
	res, err := s.Exec("insert into articles (title, uid, description, content, link, country, location, lang, source, pub_date, saved_at) values (?,?,?,?,?,?,?,?,?,?,?)",
		a.Title, a.UID, a.Description, a.Content, a.Link, a.Country, a.Location, a.Lang, a.Source, a.PubDate, a.SavedAt)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("cannot get last inserted ID for articles: %w", err)
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
	row := s.QueryRow("select a.id, a.uid, a.title, a.description, a.content, a.country, a.location, a.lang, a.source, a.pub_date from articles a where a.uid=?", uid)
	err := row.Scan(
		&article.ID,
		&article.UID,
		&article.Title,
		&article.Description,
		&article.Content,
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
	err = tx.QueryRow("select count(id) from feed_lock").Scan(&countID)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if countID == 0 {
		_, err = tx.Exec("insert into feed_lock (is_locked, timestamp) values(?,?)", true, time.Now().UnixMilli())
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		return &Lock{
			IsLocked:  false,
			Timestamp: time.Now().Add(-2 * time.Hour).UnixMilli(),
		}, tx.Commit()
	}
	res, err := tx.Exec("update feed_lock set is_locked = true where is_locked = false")

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	count, err := res.RowsAffected()

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if count == 0 {
		_ = tx.Commit()
		return nil, errors.New("already locked")
	}

	var ts int64
	err = tx.QueryRow("select timestamp from feed_lock where is_locked=true limit 1").Scan(&ts)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec("update feed_lock set is_locked = false, timestamp = ? where is_locked = true",
		time.Now().UnixMilli())
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
