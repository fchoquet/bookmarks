package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fchoquet/bookmarks/app/handlers"
	"github.com/fchoquet/bookmarks/app/middlewares"
	"github.com/fchoquet/bookmarks/app/response"
	"github.com/gorilla/mux"
)

// Start is the main application entry point.
// It could be the main() function in main.go but this approach eases injection of configuration
func Start(cfg Configuration) {

	logger = initLogger(cfg)

	logger.Info("application is starting...")

	server := &http.Server{Addr: ":8080", Handler: HTTPHandler(cfg)}

	// handles graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(err.Error())
		}
	}()

	gracefulShutdown(server, 10*time.Second)
}

// HTTPHandler returns the top level HttpHandler including all the middlewares
func HTTPHandler(cfg Configuration) http.Handler {
	sessionStore := initSessionStore()
	csrfProtection := initCSRFProtection(cfg)
	db := initDB(cfg.DBConfig)
	bookmarksRepo := initBookmarksRepo(db)
	oembedFetcher := initOembedFetcher(logger)

	r := mux.NewRouter()

	// This is shared by all the endpoints, including status
	defaultPipeline := middlewares.Pipe(
		middlewares.Timestamp,
		middlewares.Recovery,
		middlewares.RouteName,
		middlewares.TransactionID,
		// more middlewares needed (metrics, etc)
		middlewares.Log(logger),
	)

	// This is the pipeline used by the api
	apiPipeline := middlewares.Pipe(
		defaultPipeline,
		middlewares.BasicAuth(cfg.BasicAuthUsers),
	)

	// This is the pipeline used by the web interface
	webPipeline := middlewares.Pipe(
		defaultPipeline,
		middlewares.Session(sessionStore),
		csrfProtection,
	)

	r.Handle("/healthcheck",
		defaultPipeline(handlers.GetHealthcheck())).
		Methods("GET").
		Name("get_healthcheck")

	// API
	r.Handle("/bookmarks",
		apiPipeline(handlers.ListBookmarks(bookmarksRepo))).
		Methods("GET").
		Name("get_bookmarks")

	r.Handle("/bookmarks/{id}",
		apiPipeline(handlers.GetBookmark(bookmarksRepo))).
		Methods("GET").
		Name("get_bookmark")

	r.Handle("/bookmarks",
		apiPipeline(handlers.PostBookmark(bookmarksRepo, oembedFetcher))).
		Methods("POST").
		Name("post_bookmarks")

	r.Handle("/bookmarks/{id}",
		apiPipeline(handlers.DeleteBookmark(bookmarksRepo))).
		Methods("DELETE").
		Name("delete_bookmark")

	r.Handle("/bookmarks/{id}/keywords",
		apiPipeline(handlers.PutBookmarkKeywords(bookmarksRepo))).
		Methods("PUT").
		Name("put_bookmark_keywords")

	// Web
	r.Handle("/",
		webPipeline(handlers.GetIndex())).
		Methods("GET").
		Name("get_index")

	web := r.PathPrefix("/web").Subrouter()

	web.Handle("/bookmarks",
		webPipeline(handlers.GetBookmarks(bookmarksRepo))).
		Methods("GET").
		Name("get_index")

	web.Handle("/bookmarks/new",
		webPipeline(handlers.GetNewBookmark())).
		Methods("GET").
		Name("get_bookmarks_new")

	web.Handle("/bookmarks/create",
		webPipeline(handlers.PostCreateBookmark(bookmarksRepo, oembedFetcher))).
		Methods("POST").
		Name("post_bookmarks_create")

	web.Handle("/bookmarks/{id}/edit",
		webPipeline(handlers.GetEditBookmark(bookmarksRepo))).
		Methods("GET").
		Name("get_bookmarks_edit")

	web.Handle("/bookmarks/{id}/update",
		webPipeline(handlers.PostUpdateBookmark(bookmarksRepo))).
		Methods("POST").
		Name("post_bookmarks_update")

	web.Handle("/bookmarks/{id}/delete",
		webPipeline(handlers.PostDeleteBookmark(bookmarksRepo))).
		Methods("POST").
		Name("post_bookmarks_update")

	// The /docs endpoint is a little special and does not use any middleware
	// (we don't want logs, metrics and other stuff for it)
	r.PathPrefix("/docs/").Handler(handlers.GetDocHandler()).Methods("GET")

	// fallback route - returns 404
	// The default router already returns a 404 but it does not go through the middleware pipeline
	// so metrics and logs are not written. Let's use this fallback route to avoid that
	r.Handle("/{route}", defaultPipeline(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})))

	return r
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	server.Shutdown(ctx)
	logger.Info("server has shut down")
}
