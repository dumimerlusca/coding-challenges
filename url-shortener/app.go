package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	dbpool   *pgxpool.Pool
	cfg      Config
	urlModel *UrlModel
}

func NewApp(cfg Config, dbpool *pgxpool.Pool) *App {
	return &App{
		cfg:      cfg,
		dbpool:   dbpool,
		urlModel: NewUrlModel(),
	}
}

func (app *App) ListenAndServe(ctx context.Context) error {
	server := http.Server{
		Addr:    app.cfg.Http.Addr,
		Handler: app.routes(),
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err := server.Shutdown(ctxTimeout)
		if err != nil {
			fmt.Println("failed to gracefully shutdown the http server", err)
		}
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	wg.Wait()

	return nil
}
