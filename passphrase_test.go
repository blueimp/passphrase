package passphrase

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	istrings "github.com/blueimp/passphrase/internal/strings"
)

func TestWrite(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i > -10; i-- {
		Write(&buffer, i)
		str := string(buffer.Bytes())
		buffer.Reset()
		if str != "" {
			t.Errorf("Expected empty passphrase, got: %s", str)
		}
	}
	for i := 1; i <= 10; i++ {
		Write(&buffer, i)
		str := string(buffer.Bytes())
		buffer.Reset()
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

func benchmarkWrite(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Write(ioutil.Discard, i)
	}
}

func benchmarkString(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		String(i)
	}
}

func BenchmarkWrite4(b *testing.B)    { benchmarkWrite(4, b) }
func BenchmarkWrite16(b *testing.B)   { benchmarkWrite(16, b) }
func BenchmarkWrite64(b *testing.B)   { benchmarkWrite(64, b) }
func BenchmarkWrite256(b *testing.B)  { benchmarkWrite(256, b) }
func BenchmarkWrite1024(b *testing.B) { benchmarkWrite(1024, b) }

func BenchmarkString4(b *testing.B)    { benchmarkString(4, b) }
func BenchmarkString16(b *testing.B)   { benchmarkString(16, b) }
func BenchmarkString64(b *testing.B)   { benchmarkString(64, b) }
func BenchmarkString256(b *testing.B)  { benchmarkString(256, b) }
func BenchmarkString1024(b *testing.B) { benchmarkString(1024, b) }
