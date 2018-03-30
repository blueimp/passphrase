package main

import (
	"fmt"
	"net/http"

	"github.com/blueimp/passphrase"
	"github.com/blueimp/passphrase/internal/parse"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
		log.Errorf(appengine.NewContext(request), fmt.Sprint(err))
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}
