package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/t2-invert-index-search-Valynok/invertindex"
	"github.com/t2-invert-index-search-Valynok/request"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
)

var mainIndex invertindex.IndexType
var fileNames []string

func indexFile(directoryPath string, fileName string) invertindex.IndexType {

	fileContent, err := ioutil.ReadFile(directoryPath + "/" + fileName)

	if err != nil {
		fmt.Println(err)
	}

	text := string(fileContent)
	fmt.Println(invertindex.GetIndex(text, fileName))
	return invertindex.GetIndex(text, fileName)
}

func getFileNames(directoryRelativePath string) []string {

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

func indexFiles(filesDirectoryPath string, fileNames []string) invertindex.IndexType {
	filesIndex := make(invertindex.IndexType)
	fileIndexChannel := make(chan invertindex.IndexType, len(fileNames))
	for f := 0; f < len(fileNames); f++ {
		go func(fileName string) {
			fileIndex := indexFile(filesDirectoryPath, fileName)
			fmt.Println("found ", len(fileIndex), "words in ", fileName)
			fileIndexChannel <- fileIndex

		}(fileNames[f])
		fmt.Println("go routine for ", fileNames[f], " started")
	}

	for f := 0; f < len(fileNames); f++ {
		fileRes := <-fileIndexChannel
		fmt.Println("got from pipe", len(fileRes))
		invertindex.MergeIndex(filesIndex, fileRes)
	}

	return filesIndex
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	searchText, ok := r.URL.Query()["search"]

	if !ok {
		http.Error(w, "Could not find parameter 'search'", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	searchWords := strings.Split(searchText[0], " ")

	result := request.GetResult(searchWords, mainIndex, fileNames)
	for _, wordResult := range result {
		if wordResult.Value != 0 {
			fmt.Fprintln(w, wordResult.Filename, "; matches - ", wordResult.Value)
		}
	}
}

func main() {
	filesDirectoryPath := os.Args[1]
	fileNames = getFileNames(filesDirectoryPath)

	mainIndex = indexFiles(filesDirectoryPath, fileNames)

	http.HandleFunc("/", indexHandler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
