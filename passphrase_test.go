package passphrase

import (
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

func generatePassphrase(number int) string {
	pass, err := Passphrase(number)
	if err != nil {
		panic(err)
	}
	return pass
}

func TestPassphrase(t *testing.T) {
	for i := 0; i > -10; i-- {
		str := generatePassphrase(i)
		if str != "" {
			t.Errorf("Expected empty passphrase, got: %s", str)
		}
	}
	for i := 1; i <= 10; i++ {
		str := generatePassphrase(i)
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
