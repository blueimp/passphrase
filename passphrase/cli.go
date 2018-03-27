package main

import (
	"fmt"
	"os"

	"github.com/blueimp/passphrase"
)

const defaultNumber = 4
const maxNumber = 100000

var exit = os.Exit

func main() {
	var number int
	if len(os.Args) > 1 {
		arg := os.Args[1]
		number = passphrase.ParseNumber(
			arg,
			defaultNumber,
			maxNumber,
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
