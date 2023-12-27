package main

import (
	"context"
	"database/sql"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func (app *application) generateShortKey(ctx context.Context, url string) string {
	_, span := app.tracer.Start(ctx, "generteShortKey")
	defer span.End()
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
