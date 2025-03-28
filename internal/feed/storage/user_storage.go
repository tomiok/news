package storage

import (
	"database/sql"
	"time"
)

type UserStorage struct {
	DB *sql.DB
}

func (u *UserStorage) SaveSource(sourceURL, place string, userID int) error {
	query := `INSERT INTO users_sites 
			(site_id, url, category, has_content, country, location, user_id, enabled)
			VALUES
			(0, $1, 'actualidad', true, $2, $3, $4, true);`

	_, err := u.DB.Exec(query, sourceURL, place, place, userID)

	return err
}

func (u *UserStorage) Validate(token string) error {
	query := "select id, email, token from users u where u.token = $1;"

	row := u.DB.QueryRow(query, token)
	return row.Scan()
}

func (u *UserStorage) Save(email, token string) (string, error) {
	query := "insert into users (email, token, created_at) VALUES ($1, $2, $3);"

	_, err := u.DB.Exec(query, email, token, time.Now().UnixMilli())
	if err != nil {
		return "", err
	}

	return token, nil
}
