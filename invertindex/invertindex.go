package invertindex

import (
	"strings"

	"github.com/reiver/go-porterstemmer"
)

type IndexType map[string]map[string]int

func GetIndex(text string, fileName string) IndexType {
	index := make(IndexType)
	words := strings.Fields(text) //выделение слов и удаление знаков препинания
	for i := 0; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
		words[i] = strings.TrimSpace(words[i])
		words[i] = strings.TrimFunc(words[i], func(r rune) bool {
			return ((r >= 0 && r <= 64) || (r >= 91 && r <= 96) || (r >= 123))
		})
		words[i] = porterstemmer.StemString(words[i])
		if words[i] != "" {
			index.AddWordToIndex(words[i], fileName)
		}
	}

	return index
}

func (index IndexType) AddWordToIndex(word string, fileName string) IndexType {
	fileCounter, isExist := index[word]
	if isExist {
		fileCounter[fileName]++
	} else {
		index[word] = map[string]int{fileName: 1}
	}

	return index
}

func (main IndexType) MergeIndex(slave IndexType) IndexType {
	for word, fileCounter := range slave {
		for fileName, counter := range fileCounter {
			x, isExist := main[word]
			if isExist {
				x[fileName] += counter
			} else {
				main[word] = make(map[string]int)
				main[word][fileName] = counter
			}
		}
	}
	return main
}

func (index IndexType) FindIndex(words []string, fileNames []string) map[string]int {
	result := make(map[string]int)

	for _, word := range words {
		for _, fn := range fileNames {

			result[fn] = index[word][fn] + result[fn]

		}
	}

	return result
}
