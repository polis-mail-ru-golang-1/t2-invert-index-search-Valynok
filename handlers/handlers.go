package handlers

import (
	"net/http"
	"strings"

	"github.com/t2-invert-index-search-Valynok/model"

	"github.com/t2-invert-index-search-Valynok/invertindex"
	"github.com/t2-invert-index-search-Valynok/view"
	"go.uber.org/zap"
)

type Controller struct {
	view      view.View
	mainIndex invertindex.IndexType
	fileNames []string
	logger    *zap.SugaredLogger
	model     model.Model
}

func New(v view.View, m model.Model, index invertindex.IndexType, fileNames []string, l *zap.SugaredLogger) Controller {
	return Controller{view: v, model: m, mainIndex: index, fileNames: fileNames, logger: l}
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("121312")
	defer r.Body.Close()

	viewData := make([]view.SearchResult, 0)
	queryData := r.URL.Query()["text"]

	if len(queryData) == 0 || len(queryData[0]) == 0 {
		c.logger.Errorf("wrong query data: %v", queryData)
		c.view.ResultsView(viewData, w, "")
		return
	}

	searchText := queryData[0]
	c.logger.Infof("Got GET request with next request: %s", searchText)

	words := strings.Fields(searchText)
	normalizedWords := make([]string, len(words))
	for _, s := range words {
		normalizedWords = append(normalizedWords, invertindex.NormalizeWord(s))
	}

	result := c.model.GetCountersResult(normalizedWords)

	for _, fileresult := range result {
		if fileresult.Counter != 0 {
			viewData = append(viewData, view.SearchResult{FileName: fileresult.File, Counter: fileresult.Counter})
		}
	}
	c.view.ResultsView(viewData, w, searchText)
}

func (c Controller) UploadFileTextHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileName := r.Form["filename"][0]
	fileText := r.Form["book"][0]

	c.logger.Infof("Got POST request for %s with %d text length", fileName, len(fileText))

	index := invertindex.GetIndex(fileText, fileName)

	c.logger.Debugf("Found %d words in file %s", len(index), fileName)

	file := c.model.GetOrAddFile(fileName)

	var words []string
	for k, _ := range index {
		words = append(words, k)
	}

	c.logger.Debugf("words length is %d", len(words))

	existedWords := c.model.GetWords(words)

	c.logger.Debugf("Already known words count is %d", len(*existedWords))

	existedWordsIds := make(map[string]int)

	for _, val := range *existedWords {
		existedWordsIds[val.Word] = val.Id
	}

	wordsToCreate := make([]string, 0, len(words)-len(*existedWords))
	for _, val := range words {
		if existedWordsIds[val] == 0 {
			wordsToCreate = append(wordsToCreate, val)
		}
	}

	c.logger.Debugf("words to create count is %d", len(wordsToCreate))

	createdWords := c.model.AddWordBulk(wordsToCreate)

	c.logger.Debugf("created words count is %d", len(createdWords))

	for _, val := range createdWords {
		existedWordsIds[val.Word] = val.Id
	}

	c.logger.Debugf("all words count is %d", len(existedWordsIds))

	counters := make([]model.Counters, 0, len(index))
	for k, v := range index {
		//Logger.Debug(k)
		wordId := existedWordsIds[k]
		//Logger.Debug(wordId)
		//c.model.AddCounters(wordId, file.Id, )
		counters = append(counters, model.Counters{FileId: file.Id, WordId: wordId, Counter: v[file.Name]})
	}

	c.model.AddCountersBulk(counters)

	http.Redirect(w, r, "/", 302)
}

func (c Controller) UploadIndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c.logger.Info("Upload index request")

	c.view.UploadPage.ExecuteTemplate(w, "UploadPage", nil)
}

func (c Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	c.logger.Info("Index request")

	c.view.SearchView(w)
}
