package passphrase

import (
	"strings"
	"testing"

	istrings "github.com/blueimp/passphrase/internal/strings"
)

func TestString(t *testing.T) {
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
			if !istrings.InSlice(word, Words[:]) {
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
