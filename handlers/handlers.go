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

type Controller struct {
	view      view.View
	mainIndex invertindex.IndexType
	fileNames []string
	logger    *zap.SugaredLogger
}

func New(v view.View, index invertindex.IndexType, fileNames []string, l *zap.SugaredLogger) Controller {
	return Controller{view: v, mainIndex: index, fileNames: fileNames, logger: l}
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	viewData := make([]view.SearchResult, 0)
	queryData := r.URL.Query()["text"]
	if len(queryData) == 0 || len(queryData[0]) == 0 {
		c.view.ResultsView(viewData, w, "")
		return
	}

	searchText := queryData[0]
	c.logger.Infof("Got GET request with next request: %s", searchText)

	searchWords := strings.Split(searchText, " ")

	result := GetResult(searchWords, c.mainIndex, c.fileNames)

	for _, wordResult := range result {
		if wordResult.Value != 0 {
			viewData = append(viewData, view.SearchResult{FileName: wordResult.Filename, Counter: wordResult.Value})
		}
	}

	c.view.ResultsView(viewData, w, searchText)
}

func (c Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c.logger.Info("Index request")

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
