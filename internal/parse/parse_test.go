package parse

import (
	"strconv"
	"testing"
)

func TestNaturalNumber(t *testing.T) {
	for i := 0; i <= 10; i++ {
		number := NaturalNumber("", i, MaxInt)
		if number != i {
			t.Errorf(
				"Failed to handle positive default number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	for i := -10; i < 0; i++ {
		number := NaturalNumber(strconv.Itoa(i), 100, MaxInt)
		if number != 0 {
			t.Errorf(
				"Failed to handle negative number as parameter, got: %d, expected: %d.",
				number,
				0,
			)
		}
	}
	for i := 0; i <= 10; i++ {
		number := NaturalNumber(strconv.Itoa(i), 100, MaxInt)
		if number != i {
			t.Errorf(
				"Failed to handle positive number as parameter, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
	number := NaturalNumber(strconv.Itoa(MaxInt)+"0", 100, MaxInt)
	if number != MaxInt {
		t.Errorf(
			"Failed to handle int overflow, got: %d, expected: %d.",
			number,
			MaxInt,
		)
	}
	number = NaturalNumber(strconv.Itoa(-MaxInt-1)+"0", 100, MaxInt)
	if number != 0 {
		t.Errorf(
			"Failed to handle int underflow, got: %d, expected: %d.",
			number,
			0,
		)
	}
	number = NaturalNumber(strconv.Itoa(MaxInt), 100, MaxInt)
	if number != MaxInt {
		t.Errorf(
			"Failed to handle max int, got: %d, expected: %d.",
			number,
			MaxInt,
		)
	}
	number = NaturalNumber("banana", 100, MaxInt)
	if number != 0 {
		t.Errorf(
			"Failed to handle non int string, got: %d, expected: %d.",
			number,
			0,
		)
	}
	for i := 0; i <= 10; i++ {
		number := NaturalNumber("100", 0, i)
		if number != i {
			t.Errorf(
				"Failed to respect max number, got: %d, expected: %d.",
				number,
				i,
			)
		}
	}
}
