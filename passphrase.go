package passphrase

//go:generate go run generate.go

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
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

// Write writes a passphrase with the given number of words:
func Write(writer io.Writer, numberOfWords int) (n int, err error) {
	length := int64(len(Words))
	results := make(chan result)
	for i := 0; i < numberOfWords; i++ {
		go randomNumber(length, results)
	}
	for i := 0; i < numberOfWords; i++ {
		result := <-results
		if result.Error != nil {
			return n, result.Error
		}
		str := Words[result.Number]
		if n != 0 {
			str = " " + str
		}
		bytesWritten, err := io.WriteString(writer, str)
		n += bytesWritten
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// String returns a passphrase with the given number of words:
func String(numberOfWords int) (str string, err error) {
	var buffer bytes.Buffer
	_, err = Write(&buffer, numberOfWords)
	if err != nil {
		return "", err
	}
	return string(buffer.Bytes()), nil
}
