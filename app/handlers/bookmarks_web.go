package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cbdr/ats-hal/pager"
	"github.com/cbdr/ats-sq-api/response"
	"github.com/fchoquet/bookmarks/app/context"
	"github.com/fchoquet/bookmarks/bookmarks"
	"github.com/fchoquet/bookmarks/oembed"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const itemsPerPage = 5

// GetIndex is the default handler
func GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/bookmarks", http.StatusSeeOther)
		return
	}
}

// GetBookmarks returns the bookmarks list
func GetBookmarks(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.FormValue("page"))
		if err != nil || page < 0 {
			// pages are 1-based to match the URL directly
			page = 1
		}

		pager := pager.New(page, itemsPerPage)

		bookmarks, count, err := repo.List(bookmarks.Filter{
			Pager: pager,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		lastPage := pager.PageOf(count - 1)
		// This is not elegant, but the idiomatic way in go...
		pages := []int{}
		for i := 0; i < lastPage; i++ {
			pages = append(pages, i+1)
		}

		renderTemplate(w, r, "bookmarks_index.html", map[string]interface{}{
			"bookmarks": bookmarks,
			"url":       "/web/bookmarks?page=",
			"count":     count,
			"page":      page,
			"pages":     pages,
			"lastPage":  lastPage,
		})
	}
}

// GetNewBookmark returns the bookmarks creation form
func GetNewBookmark() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, r, "bookmarks_new.html", map[string]interface{}{})
	}
}

// PostCreateBookmark creates a new bookmark
func PostCreateBookmark(repo bookmarks.Repository, fetcher oembed.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := context.Session(r.Context())

		url := r.FormValue("url")
		keywords := r.FormValue("keywords")
		logger := log.WithField("url", url)

		// Let's decorate with oembed information
		link, err := fetcher.Fetch(url)
		if err != nil {
			logger.WithError(err).Warning("an invalid URL was submitted")
			session.AddFlash(Flash{
				Level:   FlashLevelWarning,
				Title:   "Holy guacamole!",
				Message: "This URL does not seem correct",
			})
			session.Save(r, w)
			http.Redirect(w, r, "/web/bookmarks/new", http.StatusSeeOther)
			return
		}

		b := bookmarks.FromOembed(&bookmarks.Bookmark{URL: url}, link)
		for _, kw := range strings.Split(keywords, ",") {
			b.Keywords = append(b.Keywords, bookmarks.Keyword(kw))
		}

		_, err = repo.Insert(b)
		if err != nil {
			// TODO: type assertion to check if validation error or internal server
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// back to the list
		session.AddFlash(Flash{
			Level:   FlashLevelSuccess,
			Title:   "Congratulations!",
			Message: "Bookmark successfuly added",
		})
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// GetEditBookmark retruns the edit form of a bookmark
func GetEditBookmark(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// load existing bookmark
		b, err := repo.ByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if b == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		keywords := []string{}
		for _, kw := range b.Keywords {
			keywords = append(keywords, string(kw))
		}

		renderTemplate(w, r, "bookmarks_edit.html", map[string]interface{}{
			"id":       b.ID,
			"url":      b.URL,
			"keywords": strings.Join(keywords, ","),
		})
	}
}

// PostUpdateBookmark updates a bookmark (the keywords only so far)
func PostUpdateBookmark(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := context.Session(r.Context())

		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		keywords := []bookmarks.Keyword{}
		for _, kw := range strings.Split(r.FormValue("keywords"), ",") {
			keywords = append(keywords, bookmarks.Keyword(kw))
		}

		if err := repo.UpdateKeywords(id, keywords); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// back to the list
		session.AddFlash(Flash{
			Level:   FlashLevelSuccess,
			Title:   "Congratulations!",
			Message: "Bookmark successfuly updated",
		})
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// PostDeleteBookmark deletes a bookmark
func PostDeleteBookmark(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := context.Session(r.Context())

		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err := repo.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// back to the list
		session.AddFlash(Flash{
			Level:   FlashLevelSuccess,
			Title:   "Congratulations!",
			Message: "Bookmark successfuly deleted",
		})
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
