package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/blueimp/passphrase"
)

func stringInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

func param(number int) map[string]string {
	parameter := strconv.Itoa(number)
	return map[string]string{"n": parameter}
}

func passphraseRequest(params map[string]string) events.APIGatewayProxyResponse {
	response, err := Handler(&events.APIGatewayProxyRequest{
		QueryStringParameters: params,
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
	if response.Headers["content-type"] != "text/plain; charset=utf-8" {
		t.Errorf(
			"Expected content-type \"text/plain\", got: \"%s\"",
			response.Headers["content-type"],
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
		if response.StatusCode != 200 {
			t.Errorf("Expected status code 200, got: %d", response.StatusCode)
		}
		if response.Headers["content-type"] != "text/plain; charset=utf-8" {
			t.Errorf(
				"Expected content-type \"text/plain\", got: \"%s\"",
				response.Headers["content-type"],
			)
		}
		if response.Body != "" {
			t.Errorf("Expected empty passphrase, got: %s", response.Body)
		}
	}
	for i := 1; i <= maxNumber; i++ {
		response := passphraseRequest(param(i))
		if response.StatusCode != 200 {
			t.Errorf("Expected status code 200, got: %d", response.StatusCode)
		}
		if response.Headers["content-type"] != "text/plain; charset=utf-8" {
			t.Errorf(
				"Expected content-type \"text/plain\", got: \"%s\"",
				response.Headers["content-type"],
			)
		}
		words := strings.Split(response.Body, " ")
		number := len(words)
		if number != i {
			t.Errorf("Incorrect number of words, got: %d, expected: %d.", number, i)
		}
		for _, word := range words {
			if !stringInList(word, passphrase.Words[:]) {
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
		if response.StatusCode != 200 {
			t.Errorf("Expected status code 200, got: %d", response.StatusCode)
		}
		if response.Headers["content-type"] != "text/plain; charset=utf-8" {
			t.Errorf(
				"Expected content-type \"text/plain\", got: \"%s\"",
				response.Headers["content-type"],
			)
		}
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
			if !stringInList(word, passphrase.Words[:]) {
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
