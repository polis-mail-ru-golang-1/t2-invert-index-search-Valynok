package invertindex

import (
	"fmt"
	"strings"
)

func Invertindex(s []byte) map[string]int {
	str := string(s)
	clear := strings.Split(str, " ")
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
	fmt.Println(clear)

	m := make(map[string]int)
	for i := 0; i < len(clear); i++ {
		m[clear[i]]++
	}
	return m
}

func Checking(words map[string]int, word string) bool {
	_, ok := words[word]
	return ok
}
