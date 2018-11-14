package invertindex

import (
	"strings"

	"github.com/reiver/go-porterstemmer"
)

type IndexType map[string]map[string]int

func GetIndex(text string, fileName string) IndexType {
	index := make(IndexType)
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		words[i] = NormalizeWord(words[i])
		if words[i] != "" {
			index.AddWordToIndex(words[i], fileName)
		}
	}

	return index
}

func NormalizeWord(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	str = strings.TrimFunc(str, func(r rune) bool {
		return ((r >= 0 && r <= 64) || (r >= 91 && r <= 96) || (r >= 123))
	})
	str = porterstemmer.StemString(str)
	return str
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
