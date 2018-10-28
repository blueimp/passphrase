package passphrase

//go:generate go run generate.go

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"sync"
)

// MaxWorkerCount sets the concurrency limit for random number generation:
var MaxWorkerCount = 128

var wordsCount = int64(len(Words))

type result struct {
	Number int64
	Error  error
}

func generateRandomNumber(maxSize int64, results chan result) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(maxSize))
	if err != nil {
		results <- result{0, err}
		return
	}
	results <- result{bigInt.Int64(), nil}
}

func generateRandomNumbers(maxSize int64, results chan result, count int) {
	if count <= MaxWorkerCount {
		for i := 0; i < count; i++ {
			go generateRandomNumber(maxSize, results)
		}
		return
	}
	tasks := make(chan int)
	var wg sync.WaitGroup
	for worker := 0; worker < MaxWorkerCount; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range tasks {
				generateRandomNumber(maxSize, results)
			}
		}()
	}
	for i := 0; i < count; i++ {
		tasks <- i
	}
	close(tasks)
	wg.Wait()
}

// Write writes a passphrase with the given number of words:
func Write(writer io.Writer, numberOfWords int) (n int, err error) {
	results := make(chan result)
	go generateRandomNumbers(wordsCount, results, numberOfWords)
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
