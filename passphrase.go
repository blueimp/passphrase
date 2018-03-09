package passphrase

//go:generate go run generate.go

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type result struct {
	Number int64
	Error  error
}

func randomNumber(maxLength int64, results chan result) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(maxLength))
	if err != nil {
		results <- result{0, err}
	}
	results <- result{bigInt.Int64(), nil}
}

// Passphrase returns a passphrase with the given number of words:
func Passphrase(numberOfWords int) (str string, err error) {
	length := int64(len(Words))
	var selectedWords = []string{}
	results := make(chan result)
	for i := 0; i < numberOfWords; i++ {
		go randomNumber(length, results)
	}
	for i := 0; i < numberOfWords; i++ {
		result := <-results
		if result.Error != nil {
			return "", result.Error
		}
		selectedWords = append(selectedWords, Words[result.Number])
	}
	return strings.Join(selectedWords, " "), nil
}
