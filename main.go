package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"./invertindex"
)

func main() {

	files := os.Args[1:]
	//открытие всех файлов и составление обратного индекса
	maps := make(map[string](map[string]int))
	for f := 0; f < len(files); f++ {
		fl, err := ioutil.ReadFile(files[f])
		if err != nil {
			fmt.Println(err)
		}
		maps[files[f]] = invertindex.Invertindex(fl)
	}

	fmt.Println("Введите поисковый запрос:")
	scanner := bufio.NewScanner(os.Stdin) //считыванеи строки, а не слова до первого пробела
	scanner.Scan()
	request := scanner.Text()
	request = string(request)
	req := strings.Split(request, " ")
	//поиск каждого слова по каждому файлу
	result := make(map[string]int)
	for f := 0; f < len(files); f++ {
		for r := 0; r < len(req); r++ {
			if invertindex.Checking(maps[files[f]], req[r]) {
				result[files[f]] = result[files[f]] + maps[files[f]][req[r]]
			} else {
				//fmt.Println("Not found")
			}
		}
	}

	//сортировка результата по количеству совпадений
	type keyvalue struct {
		filename string
		value    int
	}

	var sorted []keyvalue
	for k, v := range result {
		sorted = append(sorted, keyvalue{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})
	//вывод на экран
	for _, r := range sorted {
		if r.value != 0 {
			fmt.Println("-", r.filename, "; совпадений -", r.value)
		}
	}
}
