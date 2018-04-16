package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/blueimp/passphrase"
	istrings "github.com/blueimp/passphrase/internal/strings"
)

func param(number int) map[string]string {
	parameter := strconv.Itoa(number)
	return map[string]string{"n": parameter}
}

func passphraseRequest(args map[string]string) events.APIGatewayProxyResponse {
	response, err := Handler(&events.APIGatewayProxyRequest{
		QueryStringParameters: args,
	})
	if err != nil {
		panic(err)
	}
	return response
}

func TestHandler(t *testing.T) {
	response := passphraseRequest(map[string]string{})
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %d", response.StatusCode)
	}
	if response.Headers["cache-control"] != "private" {
		t.Errorf(
			"Expected cache-control \"private\", got: \"%s\"",
			response.Headers["cache-control"],
		)
	}
	if response.Headers["content-type"] != "text/plain;charset=utf-8" {
		t.Errorf(
			"Expected content-type \"text/plain;charset=utf-8\", got: \"%s\"",
			response.Headers["content-type"],
		)
	}
	hsts := "max-age=31536000;includeSubDomains;preload"
	if response.Headers["strict-transport-security"] != hsts {
		t.Errorf(
			"Expected strict-transport-security \"%s\", got: \"%s\"",
			hsts,
			response.Headers["strict-transport-security"],
		)
	}
	if response.Headers["x-content-type-options"] != "nosniff" {
		t.Errorf(
			"Expected x-content-type-options \"nosniff\", got: \"%s\"",
			response.Headers["x-content-type-options"],
		)
	}
	number := len(strings.Split(response.Body, " "))
	if number != defaultNumber {
		t.Errorf(
			"Incorrect default number of words, got: %d, expected: %d.",
			defaultNumber,
			number,
		)
	}
	for i := 0; i >= -maxNumber; i-- {
		response := passphraseRequest(param(i))
		if response.Body != "" {
			t.Errorf("Expected empty passphrase, got: %s", response.Body)
		}
	}
	for i := 1; i <= maxNumber; i++ {
		response := passphraseRequest(param(i))
		words := strings.Split(response.Body, " ")
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
		response := passphraseRequest(param(i))
		words := strings.Split(response.Body, " ")
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
