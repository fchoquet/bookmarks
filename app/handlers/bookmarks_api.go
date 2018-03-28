package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fchoquet/bookmarks/app/response"
	"github.com/fchoquet/bookmarks/bookmarks"
	"github.com/fchoquet/bookmarks/oembed"
	"github.com/gorilla/mux"
)

// TODO more user friendly error messages

// ListBookmarks returns the GET /bookmaks handler
func ListBookmarks(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bs, _, err := repo.List(bookmarks.Filter{})
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(r.Context(), w, bs, http.StatusOK)
	}
}

// GetBookmark returns the GET /bookmaks/:id handler
func GetBookmark(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			response.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			response.Error(w, "id must be numeric", http.StatusBadRequest)
			return
		}

		b, err := repo.ByID(id)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if b == nil {
			response.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		response.JSON(r.Context(), w, b, http.StatusOK)
	}
}

// PostBookmark returns the POST /bookmark handler
func PostBookmark(repo bookmarks.Repository, fetcher oembed.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var b bookmarks.Bookmark
		if err2 := json.Unmarshal(body, &b); err2 != nil {
			response.Error(w, err2.Error(), http.StatusBadRequest)
			return
		}

		// Let's decorate with oembed information
		link, err := fetcher.Fetch(b.URL)
		if err != nil {
			// There might be a way to have a graceful degradation here
			// but this is out of scope

			if notFoundErr, ok := err.(*oembed.NotFoundError); ok {
				// not sure of the HTTP status to return. StatusFailedDependency sounds ok
				response.Error(w, notFoundErr.Error(), http.StatusFailedDependency)
				return
			}
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b = *(bookmarks.FromOembed(&b, link))

		newB, err := repo.Insert(&b)
		if err != nil {
			// TODO: type assertion to check if validation error or internal server
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(r.Context(), w, newB, http.StatusCreated)
	}
}

// DeleteBookmark returns the DELETE /bookmaks/:id handler
func DeleteBookmark(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			response.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			response.Error(w, "id must be numeric", http.StatusBadRequest)
			return
		}

		// First load the bookmark
		b, err := repo.ByID(id)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if b == nil {
			response.Error(w, "bookmark not found", http.StatusNotFound)
			return
		}

		// Then deletes it
		if err := repo.Delete(id); err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Returns the bookmark in the json payload
		response.JSON(r.Context(), w, b, http.StatusOK)
	}
}

// PutBookmarkKeywords returns the PUT /bookmarks/{id}/keywords handler
func PutBookmarkKeywords(repo bookmarks.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars == nil || vars["id"] == "" {
			response.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			response.Error(w, "id must be numeric", http.StatusBadRequest)
			return
		}

		// First load the bookmark
		b, err := repo.ByID(id)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if b == nil {
			response.Error(w, "bookmark not found", http.StatusNotFound)
			return
		}

		// Reads the keywords in the body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var kws []bookmarks.Keyword
		if err := json.Unmarshal(body, &kws); err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Now updates the bookmark
		if err := repo.UpdateKeywords(id, kws); err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Returns the updated bookmark in the json payload
		b.Keywords = kws
		response.JSON(r.Context(), w, b, http.StatusOK)
	}
}
