package main

import (
	"fmt"
	"os"

	"github.com/blueimp/passphrase"
	"github.com/blueimp/passphrase/internal/parse"
)

const defaultNumber = 4

var exit = os.Exit

func main() {
	var number int
	if len(os.Args) > 1 {
		arg := os.Args[1]
		number = parse.NaturalNumber(
			arg,
			defaultNumber,
			parse.MaxInt,
		)
		if number == 0 && arg != "0" {
			fmt.Fprintln(os.Stderr, "argument is not a natural number:", arg)
			exit(1)
		}
	} else {
		number = defaultNumber
	}
	_, err := passphrase.Write(os.Stdout, number)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
	fmt.Println()
}
