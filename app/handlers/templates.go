package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/fchoquet/bookmarks/app/context"
	"github.com/gorilla/csrf"
)

var tmpl *template.Template

func init() {
	fmap := template.FuncMap{
		"formatDate": formatDate,
		"sup":        sup,
	}

	// this returns an error in tests but we can safely ignore it
	var err error
	tmpl, err = template.New("main").Funcs(fmap).ParseGlob("templates/*.html")
	if err != nil {
		log.Println(err)
	}
}

// Flash represents a Flash message
type Flash struct {
	Level   string
	Title   string
	Message string
}

// Common flash levels
const (
	FlashLevelSuccess = "success"
	FlashLevelDanger  = "danger"
	FlashLevelWarning = "warning"
)

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	// automatically appends flash messages
	if _, ok := data["flashes"]; !ok {
		session, ok := context.Session(r.Context())
		if ok {
			flashes := session.Flashes()
			session.Save(r, w)
			data["flashes"] = flashes
		}
	}

	// automatically appends CSRF field
	if _, ok := data[csrf.TemplateTag]; !ok {
		data[csrf.TemplateTag] = csrf.TemplateField(r)
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		// TODO: nicer error page
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func formatDate(t time.Time) string {
	return t.Format(time.ANSIC)
}

// go templates do not have arithmetic out of the box
func sup(a, b int) bool {
	return a > b
}
