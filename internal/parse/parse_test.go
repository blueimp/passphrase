package parse_test

import (
	"strconv"
	"testing"

	"github.com/blueimp/passphrase/internal/parse"
)

func TestNaturalNumber(t *testing.T) {
	number := parse.NaturalNumber("")
	if number != 0 {
		t.Errorf(
			"Failed to handle empty parameter, got: %d, expected: %d.",
			number,
			0,
		)
	}
	number = parse.NaturalNumber(strconv.Itoa(parse.MaxInt) + "0")
	if number != parse.MaxInt {
		t.Errorf(
			"Failed to handle int overflow, got: %d, expected: %d.",
			number,
			parse.MaxInt,
		)
	}
	number = parse.NaturalNumber(strconv.Itoa(-parse.MaxInt-1) + "0")
	if number != 0 {
		t.Errorf(
			"Failed to handle int underflow, got: %d, expected: %d.",
			number,
			0,
		)
	}
	number = parse.NaturalNumber(strconv.Itoa(parse.MaxInt))
	if number != parse.MaxInt {
		t.Errorf(
			"Failed to handle max int, got: %d, expected: %d.",
			number,
			parse.MaxInt,
		)
	}
	number = parse.NaturalNumber("banana")
	if number != 0 {
		t.Errorf(
			"Failed to handle non int string, got: %d, expected: %d.",
			number,
			0,
		)
	}
	for i := 0; i <= 10; i++ {
		number := parse.NaturalNumber(strconv.Itoa(i))
		if number != i {
			t.Errorf(
				"Failed to handle positive number as parameter, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	for i := -10; i < 0; i++ {
		number := parse.NaturalNumber(strconv.Itoa(i))
		if number != 0 {
			t.Errorf(
				"Failed to handle negative number as parameter, got: %d, expected: %d.",
				number,
				0,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := parse.NaturalNumber("", i)
		if number != i {
			t.Errorf(
				"Failed to handle default number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := parse.NaturalNumber("100", 0, i)
		if number != i {
			t.Errorf(
				"Failed to respect max number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := parse.NaturalNumber("-1", 0, parse.MaxInt, i)
		if number != i {
			t.Errorf(
				"Failed to respect min number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
}
