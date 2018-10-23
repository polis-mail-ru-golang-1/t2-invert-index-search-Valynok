package mapUtils

import (
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
