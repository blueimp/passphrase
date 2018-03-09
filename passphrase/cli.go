package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blueimp/passphrase"
)

const defaultNumber = 4

func number(args []string) (number int, err error) {
	if len(args) > 1 {
		arg := args[1]
		number, err := strconv.Atoi(arg)
		if err != nil || number < 0 {
			return 0, fmt.Errorf("argument is not a natural number: %s", arg)
		}
		return number, nil
	}
	return defaultNumber, nil
}

// Passphrase returns a passphrase of random words from a word list:
func Passphrase() (pass string, err error) {
	count, err := number(os.Args)
	if err != nil {
		return "", err
	}
	pass, err = passphrase.Passphrase(count)
	if err != nil {
		return "", err
	}
	return pass, nil
}

func main() {
	pass, err := Passphrase()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(pass)
}
