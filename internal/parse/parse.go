package parse

import (
	"strconv"
)

// MaxInt is the naximum integer value on this platform.
const MaxInt = int(^uint(0) >> 1)

// NaturalNumber interprets the given string parameter as natural number.
// The given int args set default values and constraints:
//   arg[0] => default number (defaults to 0)
//   arg[1] => max mumber     (defaults to MaxInt)
//   arg[2] => min number     (default to 0)
// The default number is returned if the parameter is empty.
// The max number is returned if the parameter exceeds its size.
// Zero is returned if the interpreted string is not a natural number.
// The default, max and min numbers are assumed to be natural numbers.
func NaturalNumber(parameter string, args ...int) int {
	var number int
	var argsLength = len(args)
	if parameter == "" {
		if argsLength == 0 {
			return 0
		}
		return args[0]
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
	minNumber := 0
	if argsLength > 2 {
		minNumber = args[2]
	}
	if number < minNumber {
		return minNumber
	}
	maxNumber := MaxInt
	if argsLength > 1 {
		maxNumber = args[1]
	}
	if number > maxNumber {
		return maxNumber
	}
	return number
}
