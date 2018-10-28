package mapUtils

import (
	"os"
	"sort"
)

//сортировка результата по количеству совпадений
type Keyvalue struct {
	Filename string
	Value    int
}

func GetOrderedFiles(values map[string]int) []Keyvalue {
	var sorted []Keyvalue
	for k, v := range values {
		sorted = append(sorted, Keyvalue{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	return sorted
}

func FilterFiles(files []string, predicate func(string) bool) []string {
	result := make([]string, 0)
	for _, item := range files {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map(files []os.FileInfo, f func(os.FileInfo) string) []string {
	vsm := make([]string, len(files))
	for i, v := range files {
		vsm[i] = f(v)
	}
	return vsm
}
