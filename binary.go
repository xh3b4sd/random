package random

// Binary returns a pseudo random binary value that is either 0 or 1.
func Binary() int {
	c := make(chan int, 1)

	select {
	case c <- 0:
	case c <- 1:
	}

	return <-c
}
