package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/blueimp/passphrase"
)

func stringInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func generatePassphrase(args []string) (pass string, err error) {
	os.Args = append([]string{"noop"}, args...)
	return Passphrase()
}

func TestPassphrase(t *testing.T) {
	str, _ := generatePassphrase([]string{})
	number := len(strings.Split(str, " "))
	if number != defaultNumber {
		t.Errorf(
			"Incorrect default number of words, got: %d, expected: %d.",
			defaultNumber,
			number,
		)
	}
	str, err := generatePassphrase([]string{"test"})
	if str != "" {
		t.Errorf("Expected empty passphrase, got: %s", str)
	}
	if err == nil {
		t.Errorf("Expected error object, got nil")
	} else if err.Error() != "argument is not a natural number: test" {
		t.Errorf("Expected \"not a natural number\" error, got: \"%s\"", err)
	}
	str, err = generatePassphrase([]string{"0"})
	if str != "" {
		t.Errorf("Expected empty passphrase, got: %s", str)
	}
	if err != nil {
		t.Errorf("Expected no error object, got: \"%s\"", err)
	}
	for i := -1; i >= -10; i-- {
		str, err := generatePassphrase([]string{strconv.Itoa(i)})
		if str != "" {
			t.Errorf("Expected empty passphrase, got: %s", str)
		}
		if err == nil {
			t.Errorf("Expected error object, got nil")
		} else if err.Error() != fmt.Sprintf("argument is not a natural number: %d", i) {
			t.Errorf("Expected \"not a natural number\" error, got: \"%s\"", err)
		}
	}
	for i := 1; i <= 10; i++ {
		str, _ := generatePassphrase([]string{strconv.Itoa(i)})
		words := strings.Split(str, " ")
		number := len(words)
		if number != i {
			t.Errorf("Incorrect number of words, got: %d, expected: %d.", number, i)
		}
		for _, word := range words {
			if !stringInList(word, passphrase.Words[:]) {
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
