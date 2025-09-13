// Package mathx provides math functions.
package mathx

import "math"

// Abs returns the absolute value of x. If x is the smallest negative number,
// then the absolute value cannot be represented in int64, so the function
// returns the largest positive number instead.
func Abs(x int64) int64 {
	if x == math.MinInt64 {
		return math.MaxInt64
	}

	if x < 0 {
		return -x
	}

	return x
}
