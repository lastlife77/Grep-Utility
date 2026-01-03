package searchutil

import "strings"

func (s *Search) matchRegexp(text string) bool {
	text = s.toCase(text)
	return s.re.MatchString(text)
}

func (s *Search) matchFixString(text string) bool {
	text = s.toCase(text)
	return strings.Contains(text, s.searchWord)
}
