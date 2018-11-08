package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	porterstemmer "github.com/reiver/go-porterstemmer"
	"github.com/t2-invert-index-search-Valynok/invertindex"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
	zap "go.uber.org/zap"
)

var MainIndex invertindex.IndexType
var fileNames []string
var Logger *zap.SugaredLogger

var requestForm = []byte(`
<html>
	<body> <center>
	<form action="/search" method="get">
		Request: <input type="text" name="text">
		<input type="submit" value="Search">
	</form> </center>
	</body>
</html>
`)

var resultForm = []byte(`
<html>
	<body> <center>
	<form>
  <input type="button" value="Go back!" onclick="history.back()">
</form> </center>
	</body>
</html>
`)

func IndexDirectory(directory string) {
	fileNames = GetFileNames(directory)

	MainIndex = IndexFiles(directory)
}

func IndexFile(directoryPath string, fileName string) invertindex.IndexType {

	fileContent, err := ioutil.ReadFile(directoryPath + "/" + fileName)

	if err != nil {
		fmt.Println(err)
	}

	text := string(fileContent)
	fmt.Println(invertindex.GetIndex(text, fileName))
	return invertindex.GetIndex(text, fileName)
}

func GetFileNames(directoryRelativePath string) []string {

	files, err := ioutil.ReadDir(directoryRelativePath)
	if err != nil {
		fmt.Println(err)
	}

	fileNames := mapUtils.Map(files, func(fi os.FileInfo) string { return fi.Name() })
	fileNames = mapUtils.FilterFiles(fileNames, func(fn string) bool {
		if strings.HasSuffix(fn, ".txt") {
			return true
		} else {
			return false
		}
	})
	return fileNames
}

func IndexFiles(filesDirectoryPath string) invertindex.IndexType {
	filesIndex := make(invertindex.IndexType)
	fileIndexChannel := make(chan invertindex.IndexType, len(fileNames))
	for f := 0; f < len(fileNames); f++ {
		go func(fileName string) {
			fileIndex := IndexFile(filesDirectoryPath, fileName)
			Logger.Infof("found ", len(fileIndex), "words in ", fileName)
			fileIndexChannel <- fileIndex

		}(fileNames[f])
		Logger.Infof("go routine for ", fileNames[f], " started")
	}

	for f := 0; f < len(fileNames); f++ {
		fileRes := <-fileIndexChannel
		Logger.Infof("got from pipe", len(fileRes))
		invertindex.MergeIndex(filesIndex, fileRes)
	}

	return filesIndex
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	searchText := r.URL.Query()["text"][0]
	Logger.Infof("Got GET request with next request: %s", searchText)

	//if !ok {
	//	http.Error(w, "Could not find parameter 'search'", http.StatusBadRequest)
	//	return
	//}
	w.Header().Set("Content-Type", "text/html")

	searchWords := strings.Split(searchText, " ")

	result := GetResult(searchWords, MainIndex, fileNames)
	fmt.Fprintln(w, "<center><h1>Results</h1></center>")
	for _, wordResult := range result {
		if wordResult.Value != 0 {
			fmt.Fprintln(w, "<center>", wordResult.Filename, "; matches - ", wordResult.Value, "</center>")
		}
	}
	w.Write(resultForm)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	Logger.Infof("Index request")
	w.Write(requestForm)
	return
}

func GetResult(words []string, index invertindex.IndexType, fileNames []string) []mapUtils.Keyvalue {
	for o, re := range words {
		words[o] = porterstemmer.StemString(re)
	}

	result := invertindex.FindIndex(index, words, fileNames)

	sortedFiles := mapUtils.GetOrderedFiles(result)

	return sortedFiles
}
