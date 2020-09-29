package random

// Interface creates pseudo random numbers. The service might implement retries
// using backoff strategies and timeouts.
type Interface interface {
	// Activate determines of the given probability f is meant to activate by
	// chance upon execution. Probabilities are given within (0 1) as produced
	// by Float. Given f 0.273, another random float x with equal precision is
	// randomly generated and compared to f. In case x is within (0 f) Activate
	// returns true.
	Activate(f float64) (bool, error)
	// Float produces a new random floating point number within the range (0 1),
	// meaning both of these boundaries are exlusive. So Float never returns 0
	// or 1, but any floating point number in between. Given p is the precision
	// of the produced floating point number. Given p 3 may result in 0.273.
	Float(p int) (float64, error)
	// Max tries to create a single pseudo random number. The generated number
	// is within the range [0 max), which means that max is exclusive. Given max
	// 10 generates random numbers between 0 and 9.
	Max(max int) (int, error)
	// NMax tries to create a list of new pseudo random numbers. n represents
	// the number of pseudo random numbers in the returned list. The generated
	// numbers are within the range [0 max), which means that max is exclusive.
	// Given max 10 generates random numbers between 0 and 9.
	NMax(n, max int) ([]int, error)
}
