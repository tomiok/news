package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type UserStorage struct {
	DB *sql.DB
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
		fmt.Println(err)
		return "", err
	}

	return token, nil
}
