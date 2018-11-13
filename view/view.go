package view

import "html/template"

type View struct {
	ResultsPage *template.Template
	SearchPage  *template.Template
	UploadPage  *template.Template
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

	v.UploadPage, err = template.ParseFiles("templates/upload.html")
	if err != nil {
		return v, err
	}

	return v, nil
}
