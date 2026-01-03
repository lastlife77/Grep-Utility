// Package searchutil provides helpers for searching strings read from a file.
package searchutil

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// Search defines options for searching strings.
type Search struct {
	searchWord string

	search func(text string, i int)

	match func(text string) bool
	re    *regexp.Regexp

	preContext      int
	afterContext    int
	isPreCtx        bool
	preCtxTextBuf   []string
	preCtxStrNumBuf []string
	afterCtxCount   int

	output    func(args ...string)
	outputArr []string

	count int

	enableStringNumber bool

	toCase func(text string) string

	isInvert bool
}

// New returns a new search with default settings.
func New(searchWord string) *Search {
	s := &Search{
		searchWord: searchWord,
		toCase:     skipCase,
	}
	s.output = s.defaultOutput
	s.search = s.searchDefault

	s.match = s.matchRegexp
	var err error
	s.re, err = regexp.Compile(searchWord)
	if err != nil {
		log.Fatal("Regular expression compilation error:", err)
	}

	return s
}

// AddContext adds surrounding lines to search results.
// The pre parameter specifies how many lines before,
// and after specifies how many lines after the match to include.
func (s *Search) AddContext(pre, after int) {
	s.preContext = pre
	s.afterContext = after
	s.search = s.searchInFileWithContext
	s.isPreCtx = false
	s.preCtxTextBuf = make([]string, s.preContext)
	s.preCtxStrNumBuf = make([]string, s.preContext)
	s.afterCtxCount = 0
}

// EnableOutputToArray enables output into an array.
func (s *Search) EnableOutputToArray() {
	s.outputArr = []string{}
	s.output = s.toArrayOutput
}

// GetArrayOutput returns the output array.
func (s *Search) GetArrayOutput() []string {
	return s.outputArr
}

// EnableCountOutput enables the output to display only the count of matches.
func (s *Search) EnableCountOutput() {
	s.output = s.countOutput
}

// GetCountOutput returns the count of matches found.
func (s *Search) GetCountOutput() int {
	return s.count
}

// EnableStringNumberOutput enables the output to display number of found strings.
func (s *Search) EnableStringNumberOutput() {
	s.enableStringNumber = true
}

// IgnoreCase enables case-insensitive search.
func (s *Search) IgnoreCase() {
	s.toCase = toLowerCase
	var err error
	s.searchWord = strings.ToLower(s.searchWord)
	if s.re != nil {
		s.re, err = regexp.Compile(s.searchWord)
		if err != nil {
			log.Fatal("Regular expression compilation error:", err)
		}
	}
}

// MatchFixString allows you to treat a template as a fixed string, rather than as a regex.
func (s *Search) MatchFixString() {
	s.match = s.matchFixString
	s.re = nil
}

// Invert inverts the filter; outputs lines that do not contain a template.
func (s *Search) Invert() {
	s.isInvert = true
	if s.preContext != 0 || s.afterContext != 0 {
		s.search = s.searchInFileWithContextInvert
	} else {
		s.search = s.searchDefaultInvert
	}
}

// SearchInFile searches strings in a file.
func (s *Search) SearchInFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		text := scanner.Text()
		if s.enableStringNumber {
			s.search(text, i)
			i++
		} else {
			s.search(text, 0)
		}
	}
	if s.isInvert && s.preContext > 0 && s.afterCtxCount <= 0 {
		for i := len(s.preCtxTextBuf) - 1; i >= 0; i-- {
			if s.preCtxTextBuf[i] != "" {
				s.output(s.preCtxTextBuf[i], s.preCtxStrNumBuf[i])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("File reading error:", err)
	}
}
