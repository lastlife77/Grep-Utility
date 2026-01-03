package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"testing"
)

// BenchmarkConcurrency-12
// 247585              4846 ns/op            5481 B/op         18 allocs/op
// PASS
// ok      github.com/lastlife77/grep-go   2.683s
func BenchmarkConcurrency(t *testing.B) {
	file, err := os.Open("test.html")
	if err != nil {
		log.Fatal("File opening error:", err)
	}
	for i := 0; i < t.N; i++ {
		searchInFileConcurrency("Moby", file)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}
}

// BenchmarkWithoutConcurrency-12
// 254556              4433 ns/op            5454 B/op         17 allocs/op
// PASS
// ok      github.com/lastlife77/grep-go   1.650s
func BenchmarkWithoutConcurrency(t *testing.B) {
	file, err := os.Open("test.html")
	if err != nil {
		log.Fatal("File opening error:", err)
	}
	for i := 0; i < t.N; i++ {
		searchInFileWithoutConcurrency("Moby", file)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}
}

func searchInFileConcurrency(search string, file *os.File) {
	re, err := regexp.Compile(search)
	if err != nil {
		log.Fatal("Regular expression compilation error:", err)
	}

	s := bufio.NewScanner(file)

	var wg sync.WaitGroup

	for s.Scan() {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			if re.MatchString(text) {
				fmt.Println(text)
			}
		}(s.Text())
	}

	if err = s.Err(); err != nil {
		log.Fatal("File reading error:", err)
	}

	wg.Wait()
}

func searchInFileWithoutConcurrency(search string, file *os.File) {
	re, err := regexp.Compile(search)
	if err != nil {
		log.Fatal("Regular expression compilation error:", err)
	}

	s := bufio.NewScanner(file)

	for s.Scan() {
		if re.MatchString(s.Text()) {
			fmt.Println(s.Text())
		}
	}

	if err = s.Err(); err != nil {
		log.Fatal("File reading error:", err)
	}

}
