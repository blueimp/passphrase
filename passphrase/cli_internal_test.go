package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/blueimp/passphrase"
	istrings "github.com/blueimp/passphrase/internal/strings"
)

func generatePassphrase(args []string) (code int, out string, err string) {
	os.Args = append([]string{"noop"}, args...)
	outReader, outWriter, _ := os.Pipe()
	errReader, errWriter, _ := os.Pipe()
	originalOut := os.Stdout
	originalErr := os.Stderr
	defer func() {
		os.Stdout = originalOut
		os.Stderr = originalErr
	}()
	os.Stdout = outWriter
	os.Stderr = errWriter
	exit = func(c int) {
		code = c
	}
	func() {
		main()
		outWriter.Close()
		errWriter.Close()
	}()
	stdout, _ := ioutil.ReadAll(outReader)
	stderr, _ := ioutil.ReadAll(errReader)
	return code, string(stdout), string(stderr)
}

func TestMain(t *testing.T) {
	code, out, err := generatePassphrase([]string{})
	if code != 0 {
		t.Errorf("Unexpected status code, got %d, expected: %d.", code, 0)
	}
	if err != "" {
		t.Errorf("Unexpected error output: %s.", err)
	}
	number := len(strings.Split(out, " "))
	if number != defaultNumber {
		t.Errorf(
			"Incorrect default number of words, got: %d, expected: %d.",
			defaultNumber,
			number,
		)
	}
	code, out, err = generatePassphrase([]string{"test"})
	if code != 1 {
		t.Errorf("Unexpected status code, got %d, expected: %d.", code, -1)
	}
	if err != "argument is not a natural number: test\n" {
		t.Errorf("Expected \"not a natural number\" error, got: \"%s\"", err)
	}
	if out != "\n" {
		t.Errorf("Expected empty passphrase, got: %s", out)
	}
	code, out, err = generatePassphrase([]string{"0"})
	if code != 0 {
		t.Errorf("Unexpected status code, got %d, expected: %d.", code, 0)
	}
	if err != "" {
		t.Errorf("Unexpected error output: %s.", err)
	}
	if out != "\n" {
		t.Errorf("Expected empty passphrase, got: %s", out)
	}
	for i := -1; i >= -10; i-- {
		code, out, err := generatePassphrase([]string{strconv.Itoa(i)})
		if code != 1 {
			t.Errorf("Unexpected status code, got %d, expected: %d.", code, -1)
		}
		if err != fmt.Sprintf("argument is not a natural number: %d\n", i) {
			t.Errorf("Expected \"not a natural number\" error, got: \"%s\"", err)
		}
		if out != "\n" {
			t.Errorf("Expected empty passphrase, got: %s", out)
		}
	}
	for i := 1; i <= 10; i++ {
		code, out, err := generatePassphrase([]string{strconv.Itoa(i)})
		if code != 0 {
			t.Errorf("Unexpected status code, got %d, expected: %d.", code, 0)
		}
		if err != "" {
			t.Errorf("Unexpected error output: %s.", err)
		}
		words := strings.Split(strings.TrimSuffix(out, "\n"), " ")
		number := len(words)
		if number != i {
			t.Errorf("Incorrect number of words, got: %d, expected: %d.", number, i)
		}
		for _, word := range words {
			if !istrings.InSlice(word, passphrase.Words[:]) {
				t.Errorf("Passphrase word is not in the word list: %s", word)
			}
			if len(word) < passphrase.MinWordLength {
				t.Errorf(
					"Passphrase word is shorter than %d characters: %s",
					passphrase.MinWordLength,
					word,
				)
			}
		}
	}
}
