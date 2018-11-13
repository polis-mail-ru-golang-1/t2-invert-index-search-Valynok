package view

import (
	"html/template"
	"io"
)

type View struct {
	ResultsPage *template.Template
	SearchPage  *template.Template
}

type SearchResult struct {
	FileName string
	Counter  int
}

func New() (View, error) {
	v := View{}
	var err error

	v.SearchPage, err = template.ParseFiles("templates/search.html")
	if err != nil {
		return v, err
	}

	v.ResultsPage, err = template.ParseFiles("templates/results.html")
	if err != nil {
		return v, err
	}

	return v, nil
}

func (v View) ResultsView(data []SearchResult, w io.Writer, s string) {
	v.ResultsPage.ExecuteTemplate(w, "ResultsPage",
		struct {
			Title   string
			Results []SearchResult
			Request string
		}{
			Title:   "Results",
			Results: data,
			Request: s,
		})
}

func (v View) SearchView(w io.Writer) {
	v.SearchPage.ExecuteTemplate(w, "SearchPage", nil)
}
