package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type URL struct {
	ShortKey string
	LongURL  string
}

type URLModelInterface interface {
	Insert(string, string) error
	Get(string) (string, error)
}

type URLModel struct {
	DB *sql.DB
}

func (m *URLModel) Insert(ShortKey, LongURL string) error {
	stmt := `INSERT INTO urls (shortKey, longURL) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, ShortKey, LongURL)
	return err
}

func (m *URLModel) Get(shortKey string) (string, error) {
	var longURL string
	stmt := `SELECT longURL FROM urls WHERE shortKey=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, shortKey).Scan(&longURL)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}
	return longURL, nil
}
