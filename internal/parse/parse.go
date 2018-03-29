package parse

import (
	"strconv"
)

// MaxInt is the naximum integer value on this platform:
const MaxInt = int(^uint(0) >> 1)

// NaturalNumber interprets the given string parameter as natural number.
// The defaultNumber is returned if the parameter is empty.
// The maxNumber is returned if the parameter exceeds its size.
// Zero is returned if the interpreted string is not a natural number.
// The defaultNumber and maxNumber values are assumed to be natural numbers.
func NaturalNumber(parameter string, defaultNumber int, maxNumber int) int {
	var number int
	if parameter == "" {
		return defaultNumber
	}
	number, err := strconv.Atoi(parameter)
	if err != nil {
		numError, ok := err.(*strconv.NumError)
		if ok && numError.Err == strconv.ErrRange && string(parameter[0]) != "-" {
			number = MaxInt
		} else {
			number = 0
		}
	}
	if number < 1 {
		return 0
	}
	if number > maxNumber {
		return maxNumber
	}
	return number
}
