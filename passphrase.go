package passphrase

//go:generate go run generate.go

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"strconv"
)

// MaxInt is the naximum integer value on this platform:
const MaxInt = int(^uint(0) >> 1)

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

// ParseNumber interprets the given string parameter as natural number.
// The defaultNumber is returned if the parameter is empty.
// The maxNumber is returned if the parameter exceeds its size.
// Zero is returned if the interpreted string is not a natural number.
func ParseNumber(parameter string, defaultNumber int, maxNumber int) int {
	var number int
	if parameter == "" {
		number = defaultNumber
	} else {
		var err error
		number, err = strconv.Atoi(parameter)
		if err != nil {
			numError, ok := err.(*strconv.NumError)
			if ok && numError.Err == strconv.ErrRange && string(parameter[0]) != "-" {
				number = MaxInt
			} else {
				number = 0
			}
		}
	}
	if number < 1 {
		return 0
	}
	if number > maxNumber {
		return maxNumber
	}
	return number
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
