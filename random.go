package random

import (
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Budget is an error budget implementation that decides when to retry some
	// failed attempt of executing some operation. The given budget
	// implementation is guaranteed to be reusable.
	Budget budget.Interface
	// RandFunc represents a service returning random values. Here e.g.
	// crypto/rand.Int can be used.
	RandFunc func(rand io.Reader, max *big.Int) (n *big.Int, err error)
	// RandReader represents an instance of a cryptographically strong
	// pseudo-random generator. Here e.g. crypto/rand.Reader can be used.
	RandReader io.Reader

	// Timeout represents the deadline being waited during random number
	// creation before returning a timeout error.
	Timeout time.Duration
}

type Random struct {
	budget     budget.Interface
	randFunc   func(rand io.Reader, max *big.Int) (n *big.Int, err error)
	randReader io.Reader

	timeout time.Duration
}

func New(config Config) (*Random, error) {
	if config.Budget == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Budget must not be empty", config.Budget)
	}
	if config.RandFunc == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.RandFunc must not be empty", config.RandFunc)
	}
	if config.RandReader == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.RandReader must not be empty", config.RandReader)
	}

	if config.Timeout < 0 {
		return nil, tracer.Maskf(invalidConfigError, "%T.Timeout must not be negative", config.Timeout)
	}

	r := &Random{
		budget:   config.Budget,
		randFunc: config.RandFunc,

		randReader: config.RandReader,
		timeout:    config.Timeout,
	}

	return r, nil
}

func (r *Random) Activate(f float64) (bool, error) {
	x, err := r.Float(precisionFromFloat(f))
	if err != nil {
		return false, tracer.Mask(err)
	}

	return x < f, nil
}

func (r *Random) Float(p int) (float64, error) {
	var err error

	var l []int
	for retries := 0; retries < 5; retries++ {
		l, err = r.NMax(p, 10)
		if err != nil {
			return 0, tracer.Mask(err)
		}

		if !allZero(l) {
			return intsToFloat(l), nil
		}
	}

	return 0, tracer.Maskf(executionFailedError, "could not generate random float after 5 retries")
}

func (r *Random) Max(max int) (int, error) {
	var result int

	o := func() error {
		done := make(chan struct{}, 1)
		fail := make(chan error, 1)

		go func() {
			m := big.NewInt(int64(max))
			j, err := r.randFunc(r.randReader, m)
			if err != nil {
				fail <- tracer.Mask(err)
				return
			}

			result = int(j.Int64())

			done <- struct{}{}
		}()

		select {
		case <-time.After(r.timeout):
			return tracer.Maskf(timeoutError, "after %s", r.timeout)
		case err := <-fail:
			return tracer.Mask(err)
		case <-done:
			return nil
		}
	}

	err := r.budget.Execute(o)
	if err != nil {
		return 0, tracer.Mask(err)
	}

	return result, nil
}

func (r *Random) NMax(n, max int) ([]int, error) {
	var result []int

	for i := 0; i < n; i++ {
		j, err := r.Max(max)
		if err != nil {
			return nil, tracer.Mask(err)
		}

		result = append(result, j)
	}

	return result, nil
}

func allZero(l []int) bool {
	for _, i := range l {
		if i != 0 {
			return false
		}
	}

	return true
}

func intsToFloat(l []int) float64 {
	s := "0."

	for _, i := range l {
		s += strconv.Itoa(i)
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func precisionFromFloat(f float64) int {
	s := strconv.FormatFloat(f, 'f', -1, 64)

	if strings.Contains(s, ".") {
		return len(s) - 2
	}

	return 0
}
