package searchutil

import "fmt"

func (s *Search) defaultOutput(args ...string) {
	if s.enableStringNumber {
		fmt.Printf("%v: %v", args[1], args[0])
	} else {
		fmt.Println(args[0])
	}
}

func (s *Search) toArrayOutput(args ...string) {
	if s.enableStringNumber {
		s.outputArr = append(s.outputArr, fmt.Sprintf("%v: %v", args[1], args[0]))
	} else {
		s.outputArr = append(s.outputArr, args[0])
	}
}

func (s *Search) countOutput(_ ...string) {
	s.count++
}
