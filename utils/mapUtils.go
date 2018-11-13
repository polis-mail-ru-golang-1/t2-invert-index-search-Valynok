package mapUtils

import (
	"os"
)

type Keyvalue struct {
	Filename string
	Value    int
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

func Find(n int, f func(int) bool) int {
	for i := 0; i < n; i++ {
		if f(i) {
			return i
		}
	}

	return n
}
