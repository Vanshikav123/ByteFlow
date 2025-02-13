package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Vanshikav123/gosnippet.git/internal/models"
	"github.com/Vanshikav123/gosnippet.git/ui"
	"github.com/justinas/nosurf"
)

/*
templateData struct is used to hold the data that will be passed to the HTML templates.
It represents the dynamic content and other variables that will be rendered in the templates
*/
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

/*
Scenario
Imagine you have a template html/pages/home.html with this content:

<p>Current Date: {{ .CurrentYear }}</p>
<p>Formatted Date: {{ .Snippet.CreatedAt | formatDate }}</p>
Data Used in Rendering

	data := templateData{
	    CurrentYear: 2025,
	    Snippet: &models.Snippet{
	        CreatedAt: time.Date(2025, 2, 13, 14, 30, 0, 0, time.UTC),
	    },
	}

Step-by-Step Execution
Step 1: Creating the Template

ts, err := template.New("home.html").Funcs(functions).ParseFS(ui.Files, patterns...)
Creates a new template named home.html.
Registers formatDate inside the template.
Loads the base, partials, and home.html.

Step 2: Template Execution
At runtime, Go processes:

<p>Current Date: {{ .CurrentYear }}</p>
<p>Formatted Date: {{ .Snippet.CreatedAt | formatDate }}</p>
Substituting values:

<p>Current Date: 2025</p>
<p>Formatted Date: 13 Feb 2025 at 14:30</p>
*/
var functions = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("02 Jan 2006 at 15:04")
	},
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}
	/*fs.Glob(fsys fs.FS, pattern string) ([]string, error)
		  fs.Glob is a function from the io/fs package that helps find files inside an embedded filesystem
		  (or any filesystem) that match a pattern.

		 filepath.Glob(pattern string) ([]string, error)
	Works for regular (non-embedded) filesystems. */
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// filepath.Base extracts the filename from a full file path.
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",

			"html/partials/*.html",
			page,
		}
		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		/*Suppose fs.Glob found:
		pages = ["html/pages/home.html", "html/pages/about.html"]
		For home.html, we do:

		name := filepath.Base("html/pages/home.html") // "home.html"
		patterns := []string{
		    "html/base.html",        // Main layout
		    "html/partials/*.html",  // Header, footer, etc.
		    "html/pages/home.html",  // The specific page
		}
		This means:
		template.New("home.html").Funcs(functions).ParseFS(ui.Files, "html/base.html", "html/partials/*.html", "html/pages/home.html")
		Parses base.html, all partials, and home.html together.
		Creates a cached template for home.html.
		*/
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		/**/
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		/*The CSRF middleware (nosurf) validates the token before processing the request.
		If the token is missing or incorrect, the request is rejected.
		nosurf.Token(r) generates a unique CSRF token per session*/
		CSRFToken: nosurf.Token(r),
	}
}
