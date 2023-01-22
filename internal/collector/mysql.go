package collector

import (
	"database/sql"
	"time"
)

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
	//s.Exec("insert into articles values")

	return nil, nil
}
