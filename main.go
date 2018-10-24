package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	porterstemmer "github.com/reiver/go-porterstemmer"
	invertindex "github.com/woopwoop/invertindex"
	mapUtils "github.com/woopwoop/utils"
)

func main() {

	files := os.Args[1:]

	index := make(map[string](map[string]int))
	for f := 0; f < len(files); f++ {
		fileContent, err := ioutil.ReadFile(files[f])

		if err != nil {
			fmt.Println(err)
		}

		text := string(fileContent)
		words := invertindex.GetWords(text)

		invertindex.AddWordsToIndex(words, index, files[f])
	}

	fmt.Println("Введите поисковый запрос:")
	scanner := bufio.NewScanner(os.Stdin) //считыванеи строки, а не слова до первого пробела
	scanner.Scan()
	request := scanner.Text()
	request = string(request)
	req := strings.Split(request, " ")
	for o, re := range req {
		req[o] = porterstemmer.StemString(re)
	}
	result := make(map[string]int)

	for o, _ := range req {
		for _, f := range files {
			result[f] = index[req[o]][f] + result[f]

		}
	}

	sortedFiles := mapUtils.GetOrderedFiles(result)
	for _, r := range sortedFiles {
		if r.Value != 0 {
			fmt.Println("-", r.Filename, "; совпадений -", r.Value)
		}
	}
}
