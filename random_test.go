package random

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"
)

func Test_Random_NMax_Error_RandFunc(t *testing.T) {
	var err error

	var r *Random
	{
		c := Config{
			Budget: budget.NewSingle(),
			RandFunc: func(randReader io.Reader, max *big.Int) (n *big.Int, err error) {
				return nil, tracer.Mask(timeoutError)
			},
			RandReader: &bytes.Buffer{},

			Timeout: 10 * time.Millisecond,
		}

		r, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	n := 5
	max := 10

	_, err = r.NMax(n, max)
	if !IsTimeout(err) {
		t.Fatal("expected", timeoutError, "got", nil)
	}
}

func Test_Random_NMax_Error_Timeout(t *testing.T) {
	var err error

	var r *Random
	{
		c := Config{
			Budget: budget.NewSingle(),
			RandFunc: func(randReader io.Reader, max *big.Int) (n *big.Int, err error) {
				time.Sleep(200 * time.Millisecond)
				return rand.Int(randReader, max)
			},
			RandReader: rand.Reader,

			Timeout: 20 * time.Millisecond,
		}

		r, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	n := 5
	max := 10

	_, err = r.NMax(n, max)
	if !IsTimeout(err) {
		t.Fatal("expected", timeoutError, "got", nil)
	}
}

// Test_Random_NMax_Random generates 100 random numbers between 0 and 9,
// meaning that each of the numbers has to be generated multiple times. The test
// checks that we generate the same numbers at least twice and that the defined
// boundaries hold.
func Test_Random_NMax_Boundaries(t *testing.T) {
	var err error

	var r *Random
	{
		c := Config{
			Budget:     budget.NewSingle(),
			RandFunc:   rand.Int,
			RandReader: rand.Reader,

			Timeout: 1 * time.Millisecond,
		}

		r, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	n := 100
	max := 10

	newRandomNumbers, err := r.NMax(n, max)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	alreadySeen := map[int]struct{}{}

	for _, r := range newRandomNumbers {
		// Ensure the boundaries hold.
		if r >= max {
			t.Fatal("r must be smaller than max")
		}

		alreadySeen[r] = struct{}{}
	}

	// Ensure the maximum amount holds.
	l := len(alreadySeen)
	if l != 10 {
		t.Fatal("expected", 10, "got", l)
	}
}

func Test_Random_precisionFromFloat(t *testing.T) {
	testCases := []struct {
		name      string
		float     float64
		precision int
	}{
		{
			name:      "case 0",
			float:     0,
			precision: 0,
		},
		{
			name:      "case 1",
			float:     0.0,
			precision: 0,
		},
		{
			name:      "case 2",
			float:     0.3,
			precision: 1,
		},
		{
			name:      "case 3",
			float:     0.01,
			precision: 2,
		},
		{
			name:      "case 4",
			float:     0.12345678,
			precision: 8,
		},
		{
			name:      "case 5",
			float:     0.87654321,
			precision: 8,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := precisionFromFloat(tc.float)

			if p != tc.precision {
				t.Fatalf("expected %#v to equal %#v", tc.precision, p)
			}
		})
	}
}
