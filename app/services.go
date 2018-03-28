package app

import (
	"encoding/gob"
	"fmt"

	"github.com/fchoquet/bookmarks/app/handlers"
	"github.com/fchoquet/bookmarks/app/middlewares"
	"github.com/fchoquet/bookmarks/bookmarks"
	"github.com/fchoquet/bookmarks/oembed"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// This file contains service initialization. This is where DI happens
// Since all these functions are called during the application boot phase they
// are allowed to panic. Errors in these functions are mostly caused by
// unrecoverable configurations issues

// The default logger
var logger log.FieldLogger

func initLogger(cfg Configuration) log.FieldLogger {
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})

	logger = log.StandardLogger()

	return logger
}

func initSessionStore() sessions.Store {
	// Register types stored in session
	gob.Register(handlers.Flash{})

	// TODO: manage secrets
	return sessions.NewCookieStore([]byte("something-very-secret"))
}

func initCSRFProtection(cfg Configuration) middlewares.Middleware {
	// note that we can't use CSRF protection over http, only https
	// so it's optional in dev
	return csrf.Protect(cfg.CSRFSecret, csrf.Secure(!cfg.DisableCSRFProtection))
}

func initDB(cfg DatabaseConfig) *sqlx.DB {
	// see https://github.com/go-sql-driver/mysql/issues/9 for the explanation of ?parseTime=true
	configuration := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)

	return sqlx.MustConnect("mysql", configuration)
}

func initBookmarksRepo(db *sqlx.DB) bookmarks.Repository {
	return bookmarks.NewRepository(db)
}

func initOembedFetcher(logger log.FieldLogger) oembed.Fetcher {
	fetcher, err := oembed.NewFetcher(logger)
	// There might be a way to have a graceful degradation here
	// but what's the point of starting this app if we can't fetch oembed props?
	if err != nil {
		panic(err)
	}
	return fetcher
}
