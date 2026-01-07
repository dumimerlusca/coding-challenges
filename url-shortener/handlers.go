package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UrlCreationBody struct {
	Url string
}

type UrlCreationResponse struct {
	Key      string `json:"key"`
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

func (app *App) handleShortUrlCreation(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		app.httpError(w, "invalid request body", http.StatusBadRequest, err)
		return
	}

	var body UrlCreationBody

	err = json.Unmarshal(bodyData, &body)
	if err != nil {
		app.httpError(w, "invalid request body", http.StatusBadRequest, err)
		return
	}

	parsedUrl, err := url.Parse(body.Url)
	if err != nil {
		app.httpError(w, fmt.Sprintf("invalid url: %s", err.Error()), http.StatusBadRequest, err)
		return
	}

	longUrl := parsedUrl.String()
	maxRetry := 10

	i := 0
	for {
		hashInput := longUrl

		if i != 0 {
			hashInput += strconv.Itoa(int(time.Now().Unix()))
		}

		i++

		key := generateHash(hashInput, 10)
		row, err := app.urlModel.Insert(r.Context(), app.dbpool, key, longUrl)

		if err != nil {
			if errors.Is(err, ErrDuplicateUrlKey) {
				if i >= maxRetry {
					app.httpError(w, "internal server error", http.StatusInternalServerError, err)
					return
				}

				continue
			}

			app.httpError(w, "internal server error", http.StatusInternalServerError, err)
			return
		}

		var res UrlCreationResponse
		res.Key = row.Key
		res.LongUrl = row.Long
		res.ShortUrl = fmt.Sprintf("%v/%v", app.cfg.Url, row.Key)

		writeJson(w, http.StatusCreated, res)
		return
	}
}

func (app *App) handleShortUrlRedirect(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimLeft(r.URL.Path, "/")

	row, err := app.urlModel.Get(r.Context(), app.dbpool, key)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			app.httpError(w, "not found", http.StatusNotFound, err)
			return
		}

		app.httpError(w, "internal server err", 500, err)
		return
	}

	parsedUrl, err := url.Parse(row.Long)
	if err != nil {
		app.httpError(w, "internal server err", 500, err)
		return
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}

	http.Redirect(w, r, parsedUrl.String(), http.StatusFound)
}

func (app *App) httpError(w http.ResponseWriter, message string, code int, err error) {
	http.Error(w, message, code)
	fmt.Println(err)
}
