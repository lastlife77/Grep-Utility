package searchutil

import (
	"fmt"
	"os"
	"slices"
	"testing"
)

func TestSearchInFileWithContext(t *testing.T) {
	tests := []struct {
		data         []byte
		search       string
		exp          []string
		preContext   int
		afterContext int
	}{
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "three",
			preContext:   2,
			afterContext: 0,
			exp:          []string{"one", "two", "three"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "t",
			preContext:   2,
			afterContext: 0,
			exp:          []string{"one", "two", "three"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   2,
			afterContext: 2,
			exp:          []string{"one", "two", "three", "four"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 1,
			exp:          []string{"two", "three"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 0,
			exp:          []string{"two"},
		},
		{
			data:         []byte("no\n" + "yes\n" + "no\n" + "yes\n" + "no\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 1,
			exp:          []string{"no", "yes", "no", "yes", "no"},
		},
		{
			data:         []byte("yes\n" + "no\n" + "no\n" + "no\n" + "no\n" + "yes\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 2,
			exp:          []string{"yes", "no", "no", "no", "yes"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.AddContext(test.preContext, test.afterContext)
			s.EnableOutputToArray()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%q\nExpected:\n%q", act, test.exp)
			}
		})
	}
}

func TestSearchInFileEnableCounts(t *testing.T) {
	tests := []struct {
		data   []byte
		search string
		exp    int
	}{
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "three",
			exp:    1,
		},
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "o",
			exp:    2,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.EnableCountOutput()
			s.SearchInFile(file)

			act := s.GetCountOutput()
			if act != test.exp {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileEnableCountsWithContext(t *testing.T) {
	tests := []struct {
		data         []byte
		search       string
		exp          int
		preContext   int
		afterContext int
	}{
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "three",
			preContext:   2,
			afterContext: 0,
			exp:          3,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "t",
			preContext:   2,
			afterContext: 0,
			exp:          3,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   2,
			afterContext: 2,
			exp:          4,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 1,
			exp:          2,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 0,
			exp:          1,
		},
		{
			data:         []byte("no\n" + "yes\n" + "no\n" + "yes\n" + "no\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 1,
			exp:          5,
		},
		{
			data:         []byte("yes\n" + "no\n" + "no\n" + "no\n" + "no\n" + "yes\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 2,
			exp:          5,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.AddContext(test.preContext, test.afterContext)
			s.EnableCountOutput()
			s.SearchInFile(file)

			act := s.GetCountOutput()
			if act != test.exp {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileEnableStringNumbers(t *testing.T) {
	tests := []struct {
		data   []byte
		search string
		exp    []string
	}{
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "three",
			exp:    []string{"3: three"},
		},
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "o",
			exp:    []string{"1: one", "2: two"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.EnableOutputToArray()
			s.EnableStringNumberOutput()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileWithIgnoreCase(t *testing.T) {
	tests := []struct {
		data   []byte
		search string
		exp    []string
	}{
		{
			data:   []byte("One\n" + "Two\n" + "three\n"),
			search: "three",
			exp:    []string{"three"},
		},
		{
			data:   []byte("one\n" + "TWO\n" + "three\n"),
			search: "o",
			exp:    []string{"one", "TWO"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.IgnoreCase()
			s.EnableOutputToArray()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileMatchFixString(t *testing.T) {
	tests := []struct {
		data   []byte
		search string
		exp    []string
	}{
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "th",
			exp:    []string{"three"},
		},
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "on+",
			exp:    []string{},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.MatchFixString()
			s.EnableOutputToArray()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileInvert(t *testing.T) {
	tests := []struct {
		data   []byte
		search string
		exp    []string
	}{
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "three",
			exp:    []string{"one", "two"},
		},
		{
			data:   []byte("one\n" + "two\n" + "three\n"),
			search: "o",
			exp:    []string{"three"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.Invert()
			s.EnableOutputToArray()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%v\nExpected:\n%v", act, test.exp)
			}
		})
	}
}

func TestSearchInFileWithContextInvert(t *testing.T) {
	tests := []struct {
		data         []byte
		search       string
		exp          []string
		preContext   int
		afterContext int
	}{
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "three",
			preContext:   2,
			afterContext: 0,
			exp:          []string{},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "t",
			preContext:   2,
			afterContext: 0,
			exp:          []string{},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   2,
			afterContext: 2,
			exp:          []string{"five"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 1,
			exp:          []string{"one", "four", "five"},
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 0,
			exp:          []string{"one", "three", "four", "five"},
		},
		{
			data:         []byte("no\n" + "yes\n" + "no\n" + "yes\n" + "no\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 1,
			exp:          []string{},
		},
		{
			data:         []byte("yes\n" + "no\n" + "no\n" + "no\n" + "no\n" + "yes\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 2,
			exp:          []string{"no"},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.AddContext(test.preContext, test.afterContext)
			s.EnableOutputToArray()
			s.Invert()
			s.SearchInFile(file)

			act := s.GetArrayOutput()
			if !slices.Equal(act, test.exp) {
				t.Fatalf("\nActual:\n%q\nExpected:\n%q", act, test.exp)
			}
		})
	}
}

func TestSearchInFileWithContextInvertEnableCounts(t *testing.T) {
	tests := []struct {
		data         []byte
		search       string
		exp          int
		preContext   int
		afterContext int
	}{
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "three",
			preContext:   2,
			afterContext: 0,
			exp:          0,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n"),
			search:       "t",
			preContext:   2,
			afterContext: 0,
			exp:          0,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   2,
			afterContext: 2,
			exp:          1,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 1,
			exp:          3,
		},
		{
			data:         []byte("one\n" + "two\n" + "three\n" + "four\n" + "five\n"),
			search:       "two",
			preContext:   0,
			afterContext: 0,
			exp:          4,
		},
		{
			data:         []byte("no\n" + "yes\n" + "no\n" + "yes\n" + "no\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 1,
			exp:          0,
		},
		{
			data:         []byte("yes\n" + "no\n" + "no\n" + "no\n" + "no\n" + "yes\n"),
			search:       "yes",
			preContext:   1,
			afterContext: 2,
			exp:          1,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Case: %v\n", i), func(t *testing.T) {
			file := createFile(t, test.data)
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})

			s := New(test.search)
			s.AddContext(test.preContext, test.afterContext)
			s.EnableCountOutput()
			s.Invert()
			s.SearchInFile(file)

			act := s.GetCountOutput()
			if act != test.exp {
				t.Fatalf("\nActual:\n%q\nExpected:\n%q", act, test.exp)
			}
		})
	}
}

func createFile(t *testing.T, data []byte) *os.File {
	file, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err.Error())
	}

	if _, err := file.Write(data); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatalf("Failed to open file: %v", err.Error())
	}

	return file
}
