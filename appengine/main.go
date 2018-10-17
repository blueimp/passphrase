package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blueimp/passphrase"
	"github.com/blueimp/passphrase/internal/parse"
)

const defaultNumber = 4
const maxNumber = 100

func indexHandler(response http.ResponseWriter, request *http.Request) {
	number := parse.NaturalNumber(
		request.FormValue("n"),
		defaultNumber,
		maxNumber,
	)
	response.Header().Set("cache-control", "private")
	response.Header().Set("content-type", "text/plain;charset=utf-8")
	response.Header().Set(
		"strict-transport-security",
		"max-age=31536000;includeSubDomains;preload",
	)
	response.Header().Set("x-content-type-options", "nosniff")
	_, err := passphrase.Write(response, number)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
