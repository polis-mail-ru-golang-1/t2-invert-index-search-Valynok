package invertindex

import (
	"strings"
)

func Invertindex(s []byte) map[string]int {
	str := string(s)
	clear := strings.Split(str, " ") //выделение слов и удаление знаков препинания
	for i := 0; i < len(clear); i++ {
		clear[i] = strings.TrimSpace(clear[i])
		clear[i] = strings.Trim(clear[i], ",")
		clear[i] = strings.Trim(clear[i], "-")
		clear[i] = strings.Trim(clear[i], ":")
		clear[i] = strings.Trim(clear[i], ";")
		if clear[i] == "" {
			clear = append(clear[:i], clear[i+1:]...)
		}
	}
	//подсчет повторений для каждого слова
	m := make(map[string]int)
	for i := 0; i < len(clear); i++ {
		m[clear[i]]++
	}
	return m
}

//проверка встретилось ли слово в данном файле
func Checking(words map[string]int, word string) bool {
	_, ok := words[word]
	return ok
}
