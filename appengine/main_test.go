package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/blueimp/passphrase"
	istrings "github.com/blueimp/passphrase/internal/strings"
)

func param(number int) string {
	return "/?n=" + strconv.Itoa(number)
}

func passphraseRequest(url string) (response http.Response, result string) {
	request := httptest.NewRequest("GET", url, nil)
	recorder := httptest.NewRecorder()
	indexHandler(recorder, request)
	response = *recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	return response, string(body)
}

func TestPassphrase(t *testing.T) {
	response, result := passphraseRequest("/")
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %d", response.StatusCode)
	}
	if response.Header.Get("cache-control") != "private" {
		t.Errorf(
			"Expected cache-control \"private\", got: \"%s\"",
			response.Header.Get("cache-control"),
		)
	}
	if response.Header.Get("content-type") != "text/plain;charset=utf-8" {
		t.Errorf(
			"Expected content-type \"text/plain;charset=utf-8\", got: \"%s\"",
			response.Header.Get("content-type"),
		)
	}
	hsts := "max-age=31536000;includeSubDomains;preload"
	if response.Header.Get("strict-transport-security") != hsts {
		t.Errorf(
			"Expected strict-transport-security \"%s\", got: \"%s\"",
			hsts,
			response.Header.Get("strict-transport-security"),
		)
	}
	if response.Header.Get("x-content-type-options") != "nosniff" {
		t.Errorf(
			"Expected x-content-type-options \"nosniff\", got: \"%s\"",
			response.Header.Get("x-content-type-options"),
		)
	}
	number := len(strings.Split(result, " "))
	if number != defaultNumber {
		t.Errorf(
			"Incorrect default number of words, got: %d, expected: %d.",
			defaultNumber,
			number,
		)
	}
	for i := 0; i >= -maxNumber; i-- {
		_, result := passphraseRequest(param(i))
		if result != "" {
			t.Errorf("Expected empty passphrase, got: %s", response.Body)
		}
	}
	for i := 1; i <= maxNumber; i++ {
		_, result := passphraseRequest(param(i))
		words := strings.Split(result, " ")
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
	for i := maxNumber + 1; i <= maxNumber+11; i++ {
		_, result := passphraseRequest(param(i))
		words := strings.Split(result, " ")
		number := len(words)
		if number != maxNumber {
			t.Errorf(
				"Incorrect number of words, got: %d, expected: %d.",
				number,
				maxNumber,
			)
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
