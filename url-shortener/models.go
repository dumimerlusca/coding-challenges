package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var ErrDuplicateUrlKey = errors.New("duplicate url key")
var ErrNoRows = errors.New("no rows found")

type UrlRow struct {
	Key  string
	Long string
}

type UrlModel struct{}

func NewUrlModel() *UrlModel {
	return &UrlModel{}
}

func (m *UrlModel) Insert(ctx context.Context, conn DatabaseConnection, key string, longurl string) (*UrlRow, error) {

	row := conn.QueryRow(ctx, "INSERT INTO url(key, long) VALUES($1,$2) RETURNING key, long", key, longurl)

	var urlRow UrlRow
	err := row.Scan(&urlRow.Key, &urlRow.Long)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, ErrDuplicateUrlKey
		}

		return nil, err
	}

	return &urlRow, nil
}

func (M *UrlModel) Get(ctx context.Context, conn DatabaseConnection, key string) (*UrlRow, error) {
	row := conn.QueryRow(ctx, "SELECT key, long FROM url WHERE key=$1", key)

	var urlRow UrlRow
	err := row.Scan(&urlRow.Key, &urlRow.Long)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, err
	}

	return &urlRow, nil
}

func (M *UrlModel) Delete(ctx context.Context, conn DatabaseConnection, key string) error {
	return nil
}
