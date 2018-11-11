package invertindex

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetIndex(t *testing.T) {
	in := "A strong, positive self-image is the best possible preparation for success."

	fileName := "file.txt"

	expectedFilesMap := IndexType{
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

	actual := GetIndex(in, fileName)

	eq := cmp.Equal(expectedFilesMap, actual)
	if !eq {
		t.Errorf("%v is not equal to expected %v", actual, expectedFilesMap)
	}
}

func TestMergeIndex(t *testing.T) {
	main, addible, expected := getData()

	actual := MergeIndex(main, addible)

	eq := cmp.Equal(expected, actual)
	if !eq {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func getData() (IndexType, IndexType, IndexType) {
	fileName := "file.txt"
	main := IndexType{
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

	addible := IndexType{
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

	expected := IndexType{
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
