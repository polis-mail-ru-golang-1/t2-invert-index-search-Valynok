package invertindex

import (
	"strings"

	"github.com/reiver/go-porterstemmer"
)

type Fileint struct {
	filename string
	counter  int
}

func GetWords(text string) []string {
	words := strings.Split(text, " ") //выделение слов и удаление знаков препинания
	for i := 0; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
		words[i] = strings.TrimSpace(words[i])
		words[i] = strings.TrimFunc(words[i], func(r rune) bool {
			return ((r >= 0 && r <= 64) || (r >= 91 && r <= 96) || (r >= 123))
		})
		words[i] = porterstemmer.StemString(words[i])
		if words[i] == "" {
			words = append(words[:i], words[i+1:]...)
		}
	}
	return words
}

func AddWordsToIndex(words []string, index map[string]map[string]int, fileName string) {

	for _, word := range words {
		filesCounters, isExist := index[word]
		if isExist {
			filesCounters[fileName]++
		} else {
			index[word] = map[string]int{fileName: 1}
		}

	}
}
