// A simple grep-like utility written in Go.
// Searches for a string in files or standard input.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/lastlife77/grep-go/internal/searchutil"
)

func main() {
	log.SetFlags(0)

	var file *os.File

	a := flag.Int("A", 0, "After each line found, additionally output N lines after it.")
	b := flag.Int("B", 0, "Output N lines to each found line.")
	ctx := flag.Int("C", 0, "Output N lines of context around the found string.")
	c := flag.Bool("c", false, "Output only the number of lines that match the pattern.")
	n := flag.Bool("n", false, "Output the line number before each found line.")
	i := flag.Bool("i", false, "Ignore the case.")
	f := flag.Bool("f", false, "Treat a template as a fixed string rather than a regular expression.")
	v := flag.Bool("v", false, "Invert the filter: output lines that do not contain a template.")

	flag.Parse()
	if *ctx != 0 && *a != 0 {
		log.Fatal("The C and A flags do not match.")
	}
	if *ctx != 0 && *b != 0 {
		log.Fatal("The C and B flags do not match.")
	}
	if *c && *n {
		log.Fatal("The c and n flags do not match.")
	}
	args := flag.Args()
	search := args[0]

	if len(args) > 1 {
		var err error
		file, err = os.Open(args[1])
		if err != nil {
			log.Fatal("File opening error:", err)
		}
	} else {
		file = os.Stdin
	}

	s := searchutil.New(search)

	if *ctx != 0 {
		s.AddContext(*ctx, *ctx)
	}
	if *a != 0 || *b != 0 {
		s.AddContext(*a, *b)
	}
	if *c {
		s.EnableCountOutput()
	}
	if *n {
		s.EnableStringNumberOutput()
	}
	if *i {
		s.IgnoreCase()
	}
	if *f {
		s.MatchFixString()
	}
	if *v {
		s.Invert()
	}

	s.SearchInFile(file)
}
