package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := Config{}

	flag.StringVar(&cfg.Http.Addr, "addr", ":8081", "addr for the http server to listen to")
	flag.StringVar(&cfg.Db.ConnString, "dsn", "postgres://local_user:local_password@localhost:5432/url_shortener", "database connection string")
	flag.StringVar(&cfg.Url, "url", "localhost:8081", "complete url")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)

	go func() {
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit
		cancel()
	}()

	dbpool, err := NewPool(ctx, cfg.Db)
	if err != nil {
		panic(err)
	}

	app := NewApp(cfg, dbpool)

	err = app.ListenAndServe(ctx)
	if err != nil {
		panic(err)
	}
}
