package storage

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"news/internal/feed"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	maxOpenConnections = 25
	maxIdleConnections = 20
)

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

func (s *SQLStorage) SaveArticle(a feed.Article) (feed.Article, error) {
	loc := strings.ToLower(a.Location)
	_, err := s.Exec(`insert into articles 
    (title, uid, description, content, raw_content, link, country, location, lang, site_id, pub_date, saved_at,categories, n_search) 
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13, to_tsvector($14))`,
		a.Title,
		a.UID,
		a.Description,
		a.Content,
		a.RawContent,
		a.Link,
		strings.ToLower(a.Country),
		loc,
		a.Lang,
		a.SourceID,
		a.PubDate,
		a.SavedAt,
		pq.Array(a.Categories),
		index(a.Categories, loc, string(a.Content)))

	if err != nil {
		return feed.Article{}, err
	}

	return a, nil
}

func (s *SQLStorage) GetArticleByUID(uid string) (feed.Article, error) {
	var article feed.Article
	row := s.QueryRow("select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.country, a.location, a.lang, a.site_id, a.pub_date, a.categories from articles a where a.uid=$1", uid)
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
		&article.SourceID,
		&article.PubDate,
		pq.Array(&categories),
	)

	if err != nil {
		return article, fmt.Errorf("cannot get articles %w", err)
	}

	article.Categories = categories
	return article, nil
}

const (
	querySelect = `select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.link, a.country, a.location, a.lang, a.pub_date, a.categories, ts_rank(n_search, query) as rank from articles a, to_tsquery('caba | rosario') query where n_search @@ query and pub_date >= $1 order by rank desc`
	rankQuery   = `select a.id, a.uid, a.title, a.description, a.content, a.raw_content, a.link, a.country, a.location, a.lang, a.pub_date, a.categories, ts_rank(n_search, query) as rank from articles a, to_tsquery($1) query where n_search  @@ query order by rank desc`
	defSize     = 50
)

func (s *SQLStorage) GetDBFeed(q string) ([]feed.Article, error) {
	back24Hours := time.Now().Add(-time.Hour * 24).UnixMilli()
	var (
		query string
		p     any
	)

	if q != "" {
		query = rankQuery
		p = formatInputQuery(q)
	} else {
		query = querySelect
		p = back24Hours
	}

	rows, err := s.Query(query, p)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var categories []string
	result := make([]feed.Article, 0, defSize)
	for rows.Next() {
		var article feed.Article
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
			new(float64),
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

func (s *SQLStorage) GetSites() ([]feed.Site, error) {
	rows, err := s.Query("select id, url, category, has_content, country, location from sites")

	if err != nil {
		return nil, err
	}

	var result []feed.Site
	for rows.Next() {
		var site feed.Site
		err = rows.Scan(&site.ID, &site.URL, &site.MainCategory, &site.HasContent, &site.Country, &site.Location)
		if err != nil {
			log.Error().Err(err)
		}
		result = append(result, site)
	}

	return result, nil
}

// formatInputQuery this is only with AND (& operator)
func formatInputQuery(q string) string {
	res := strings.Split(q, " ")

	return strings.Join(res, " & ")
}

func index(categories []string, loc, content string) string {
	var strBuilder strings.Builder

	if len(categories) > 0 {
		strBuilder.WriteString(strings.Join(categories, " "))
	}
	strBuilder.WriteString(loc)
	strBuilder.WriteString(content)

	return strBuilder.String()
}
