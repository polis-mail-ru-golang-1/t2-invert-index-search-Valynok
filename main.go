package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	porterstemmer "github.com/reiver/go-porterstemmer"
	invertindex "github.com/t2-invert-index-search-Valynok/invertindex"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
)

func indexFile(fileName string) map[string](map[string]int) {
	fileContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
	}

	text := string(fileContent)
	words := invertindex.GetWords(text)
	fmt.Println("found ", len(words), "words in file ", fileName)
	return invertindex.AddWordsToIndex(words, fileName)
}

func main() {

	files := os.Args[1:]

	index := make(map[string](map[string]int))
	fileIndexes := make(chan map[string]map[string]int, len(files))
	var wg sync.WaitGroup
	for f := 0; f < len(files); f++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, fileName string) {
			fileIndexes <- indexFile(fileName)
			defer wg.Done()
		}(&wg, files[f])
		fmt.Println("go routine for ", files[f], " started")
	}

	fmt.Println("waiting for go routines...")
	wg.Wait()

	for f := 0; f < len(files); f++ {

		fileRes := <-fileIndexes
		//fmt.Println(fileRes)
		for word, fileCounter := range fileRes {
			for fileName, counter := range fileCounter {
				x, isExist := index[word]
				if isExist {
					x[fileName] += counter
				} else {
					index[word] = make(map[string]int)
					index[word][fileName] = counter
				}

			}
			//fmt.Println(fileCounter)

			//fmt.Println(fileCounter[fileName])
		}

	}
	//fmt.Println(index)
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
