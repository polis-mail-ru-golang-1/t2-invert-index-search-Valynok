package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/t2-invert-index-search-Valynok/model"

	porterstemmer "github.com/reiver/go-porterstemmer"
	"github.com/t2-invert-index-search-Valynok/invertindex"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
	"github.com/t2-invert-index-search-Valynok/view"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

type Controller struct {
	view  view.View
	model model.Model
}

func New(v view.View, m model.Model) Controller {
	return Controller{view: v, model: m}
}

func (c Controller) ResultsView(str string, w io.Writer, s string) {
	c.view.ResultsPage.ExecuteTemplate(w, "ResultsPage",
		struct {
			Title   string
			Results string
			Request string
		}{
			Title:   "Results",
			Results: str,
			Request: s,
		})
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	searchText := r.URL.Query()["text"][0]
	Logger.Infof("Got GET request with next request: %s", searchText)

	searchWords := strings.Split(searchText, " ")

	indexedWords := c.model.GetWords(searchWords)

	// result := GetResult(searchWords, MainIndex, FileNames)
	str := ""
	// for _, wordResult := range result {
	// 	if wordResult.Value != 0 {
	// 		str += (wordResult.Filename + "; matches - " + strconv.Itoa(wordResult.Value) + "\n")
	// 	}
	// }

	c.ResultsView(str, w, searchText)

}

func (c Controller) UploadFileTextHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileName := r.Form["filename"][0]
	fileText := r.Form["book"][0]

	Logger.Infof("Got POST request for %s with %d text length", fileName, len(fileText))

	index := invertindex.GetIndex(fileText, fileName)

	Logger.Debugf("Found %d words in file %s", len(index), fileName)

	file := c.model.GetOrAddFile(fileName)

	var words []string
	for k, _ := range index {
		words = append(words, k)
	}

	Logger.Debugf("words length is %d", len(words))

	existedWords := c.model.GetWords(words)

	Logger.Debugf("Already known words count is %d", len(*existedWords))

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

	Logger.Debugf("words to create count is %d", len(wordsToCreate))

	createdWords := c.model.AddWordBulk(wordsToCreate)

	Logger.Debugf("created words count is %d", len(createdWords))

	for _, val := range createdWords {
		existedWordsIds[val.Word] = val.Id
	}

	// allWords := make([]model.Word, 0, len(words))

	// allWords = append(allWords, (*existedWords)...)
	// allWords = append(allWords, createdWords...)

	// Logger.Debug(createdWords)
	// Logger.Debug(existedWords)
	// Logger.Debug(allWords)

	Logger.Debugf("all words count is %d", len(existedWordsIds))

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

	Logger.Info("Upload index request")

	c.view.UploadPage.ExecuteTemplate(w, "UploadPage", nil)
}

func (c Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	Logger.Info("Index request")

	c.view.SearchPage.ExecuteTemplate(w, "SearchPage", nil)
}

func GetResult(words []string, index invertindex.IndexType, fileNames []string) []mapUtils.Keyvalue {
	for o, re := range words {
		words[o] = porterstemmer.StemString(re)
	}

	result := invertindex.FindIndex(index, words, fileNames)

	sortedFiles := mapUtils.GetOrderedFiles(result)

	return sortedFiles
}
