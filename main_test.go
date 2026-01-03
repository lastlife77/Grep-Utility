package main

import (
	"os"
	"os/exec"
	"slices"
	"testing"
)

func TestWithoutFlags(t *testing.T) {
	file, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("File opening error: %v", err.Error())
	}
	t.Cleanup(func() {
		os.Remove(file.Name())
	})
	data := []byte(
		`All this while Tashtego, Daggoo, and Queequeg had looked on with even more
intense interest and surprise than the rest, and at the mention of the
wrinkled brow and crooked jaw they had started as if each was separately
touched by some specific recollection.`)
	search := "th+"
	exp := []byte(
		"All this while Tashtego, Daggoo, and Queequeg had looked on with even more\n" +
			"intense interest and surprise than the rest, and at the mention of the\n" +
			"wrinkled brow and crooked jaw they had started as if each was separately\n")

	if _, err := file.Write(data); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	cmd := exec.Command("go", "run", "main.go", search, file.Name())
	act, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if !slices.Equal(act, exp) {
		t.Fatalf("\nActual:\n%q\nExpected:\n%q", act, exp)
	}
}

func TestMismatchedFlags(t *testing.T) {
	file, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatalf("File opening error: %v", err.Error())
	}

	t.Cleanup(func() {
		os.Remove(file.Name())
	})
	data := []byte(
		`All this while Tashtego, Daggoo, and Queequeg had looked on with even more
intense interest and surprise than the rest, and at the mention of the
wrinkled brow and crooked jaw they had started as if each was separately
touched by some specific recollection.`)
	search := "th+"
	exp := []byte("The C and A flags do not match.\nexit status 1\n")

	if _, err := file.Write(data); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}

	cmd := exec.Command("go", "run", "main.go", "-C=2", "-A=2", search, file.Name())
	act, err := cmd.CombinedOutput()
	if err.Error() != "exit status 1" {
		t.Fatalf("Unexpected error: %v", err.Error())
	}
	if !slices.Equal(act, exp) {
		t.Fatalf("\nActual:\n%q\nExpected:\n%q", act, exp)
	}
}
