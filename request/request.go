package request

import (
	porterstemmer "github.com/reiver/go-porterstemmer"
	"github.com/t2-invert-index-search-Valynok/invertindex"
	"github.com/t2-invert-index-search-Valynok/utils"
)

func GetResult(words []string, index invertindex.IndexType, fileNames []string) []mapUtils.Keyvalue {
	for o, re := range words {
		words[o] = porterstemmer.StemString(re)
	}

	result := invertindex.FindIndex(index, words, fileNames)

	sortedFiles := mapUtils.GetOrderedFiles(result)

	return sortedFiles
}
