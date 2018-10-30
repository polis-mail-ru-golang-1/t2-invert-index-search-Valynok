package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/t2-invert-index-search-Valynok/invertindex"
)

func TestGetIndex(t *testing.T) {
	in := "A strong, positive self-image is the best possible preparation for success."

	fileName := "file.txt"

	expectedFilesMap := invertindex.IndexType{
		"a": map[string]int{
			fileName: 1,
		},
		"possibl": map[string]int{
			fileName: 1,
		},
		"prepar": map[string]int{
			fileName: 1,
		},
		"strong": map[string]int{
			fileName: 1,
		},
		"posit": map[string]int{
			fileName: 1,
		},
		"self-imag": map[string]int{
			fileName: 1,
		},
		"is": map[string]int{
			fileName: 1,
		},
		"the": map[string]int{
			fileName: 1,
		},
		"best": map[string]int{
			fileName: 1,
		},
		"for": map[string]int{
			fileName: 1,
		},
		"success": map[string]int{
			fileName: 1,
		},
	}

	actual := invertindex.GetIndex(in, fileName)

	eq := cmp.Equal(expectedFilesMap, actual)
	if !eq {
		t.Errorf("%v is not equal to expected %v", actual, expectedFilesMap)
	}
}

func TestMergeIndex(t *testing.T) {
	main, addible, expected := getData()

	actual := invertindex.MergeIndex(main, addible)

	eq := cmp.Equal(expected, actual)
	if !eq {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func getData() (invertindex.IndexType, invertindex.IndexType, invertindex.IndexType) {
	fileName := "file.txt"
	main := invertindex.IndexType{
		"a": map[string]int{
			fileName: 2,
		},
		"possibl": map[string]int{
			fileName: 1,
		},

		"strong": map[string]int{
			fileName: 1,
		},
		"posit": map[string]int{
			fileName: 2,
		},
		"self-imag": map[string]int{
			fileName: 1,
		},
		"is": map[string]int{
			fileName: 1,
		},
		"the": map[string]int{
			fileName: 1,
		},
		"best": map[string]int{
			fileName: 1,
		},
		"for": map[string]int{
			fileName: 1,
		},
		"success": map[string]int{
			fileName: 1,
		},
	}

	addible := invertindex.IndexType{
		"a": map[string]int{
			fileName: 1,
		},
		"prepar": map[string]int{
			fileName: 1,
		},
		"strong": map[string]int{
			fileName: 1,
		},
		"posit": map[string]int{
			fileName: 3,
		},
		"self-imag": map[string]int{
			fileName: 1,
		},
		"is": map[string]int{
			fileName: 1,
		},
		"the": map[string]int{
			fileName: 1,
		},
		"best": map[string]int{
			fileName: 1,
		},
		"for": map[string]int{
			fileName: 1,
		},
		"success": map[string]int{
			fileName: 1,
		},
	}

	expected := invertindex.IndexType{
		"a": map[string]int{
			fileName: 3,
		},
		"possibl": map[string]int{
			fileName: 1,
		},
		"prepar": map[string]int{
			fileName: 1,
		},
		"strong": map[string]int{
			fileName: 2,
		},
		"posit": map[string]int{
			fileName: 5,
		},
		"self-imag": map[string]int{
			fileName: 2,
		},
		"is": map[string]int{
			fileName: 2,
		},
		"the": map[string]int{
			fileName: 2,
		},
		"best": map[string]int{
			fileName: 2,
		},
		"for": map[string]int{
			fileName: 2,
		},
		"success": map[string]int{
			fileName: 2,
		},
	}
	return main, addible, expected
}
