package handlers

import (
	"net/http"
	"strings"

	porterstemmer "github.com/reiver/go-porterstemmer"
	"github.com/t2-invert-index-search-Valynok/invertindex"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
	"github.com/t2-invert-index-search-Valynok/view"
	"go.uber.org/zap"
)

var MainIndex invertindex.IndexType
var FileNames []string
var Logger *zap.SugaredLogger

type Controller struct {
	view view.View
}

func New(v view.View) Controller {
	return Controller{view: v}
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	searchText := r.URL.Query()["text"][0]
	Logger.Infof("Got GET request with next request: %s", searchText)

	searchWords := strings.Split(searchText, " ")

	result := GetResult(searchWords, MainIndex, FileNames)

	viewData := make([]view.SearchResult, 0)
	for _, wordResult := range result {
		if wordResult.Value != 0 {
			viewData = append(viewData, view.SearchResult{FileName: wordResult.Filename, Counter: wordResult.Value})
		}
	}

	c.view.ResultsView(viewData, w, searchText)
}

func (c Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	Logger.Info("Index request")

	c.view.SearchView(w)
}

func GetResult(words []string, index invertindex.IndexType, fileNames []string) []mapUtils.Keyvalue {
	for o, re := range words {
		words[o] = porterstemmer.StemString(re)
	}

	result := invertindex.FindIndex(index, words, fileNames)

	sortedFiles := mapUtils.GetOrderedFiles(result)

	return sortedFiles
}
