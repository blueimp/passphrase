package passphrase

import (
	"strconv"
	"strings"
	"testing"
)

func stringInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func TestParseNumber(t *testing.T) {
	for i := -10; i < 0; i++ {
		number := ParseNumber("", i, MaxInt)
		if number != 0 {
			t.Errorf(
				"Failed to handle negative default number, got: %d, expected: %d.",
				number,
				0,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := ParseNumber("", i, MaxInt)
		if number != i {
			t.Errorf(
				"Failed to handle positive default number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	for i := -10; i < 0; i++ {
		number := ParseNumber(strconv.Itoa(i), 100, MaxInt)
		if number != 0 {
			t.Errorf(
				"Failed to handle negative number as parameter, got: %d, expected: %d.",
				number,
				0,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := ParseNumber(strconv.Itoa(i), 100, MaxInt)
		if number != i {
			t.Errorf(
				"Failed to handle positive number as parameter, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	number := ParseNumber(strconv.Itoa(MaxInt)+"0", 100, MaxInt)
	if number != MaxInt {
		t.Errorf(
			"Failed to handle int overflow, got: %d, expected: %d.",
			number,
			MaxInt,
		)
	}
	number = ParseNumber(strconv.Itoa(-MaxInt-1)+"0", 100, MaxInt)
	if number != 0 {
		t.Errorf(
			"Failed to handle int underflow, got: %d, expected: %d.",
			number,
			0,
		)
	}
	number = ParseNumber(strconv.Itoa(MaxInt), 100, MaxInt)
	if number != MaxInt {
		t.Errorf(
			"Failed to handle max int, got: %d, expected: %d.",
			number,
			MaxInt,
		)
	}
	number = ParseNumber("banana", 100, MaxInt)
	if number != 0 {
		t.Errorf(
			"Failed to handle non int string, got: %d, expected: %d.",
			number,
			0,
		)
	}
	number = ParseNumber("", 100, 10)
	if number != 10 {
		t.Errorf(
			"Failed to respect max number, got: %d, expected: %d.",
			number,
			10,
		)
	}
	for i := 0; i <= 10; i++ {
		number := ParseNumber("100", 0, i)
		if number != i {
			t.Errorf(
				"Failed to respect max number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
}

func TestPassphrase(t *testing.T) {
	for i := 0; i > -10; i-- {
		str, _ := String(i)
		if str != "" {
			t.Errorf("Expected empty passphrase, got: %s", str)
		}
	}
	for i := 1; i <= 10; i++ {
		str, _ := String(i)
		words := strings.Split(str, " ")
		number := len(words)
		if number != i {
			t.Errorf("Incorrect number of words, got: %d, expected: %d.", number, i)
		}
		for _, word := range words {
			if !stringInList(word, Words[:]) {
				t.Errorf("Passphrase word is not in the word list: %s", word)
			}
			if len(word) < MinWordLength {
				t.Errorf(
					"Passphrase word is shorter than %d characters: %s",
					MinWordLength,
					word,
				)
			}
		}
	}
}
