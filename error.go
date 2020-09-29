package random

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

// executionFailedError is an error type for unexpected situations where further
// code execution cannot continue. This error should never be matched against.
// Therefore there is no error matcher implemented.
var executionFailedError = &tracer.Error{
	Kind: "executionFailedError",
}

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}

var timeoutError = &tracer.Error{
	Kind: "timeoutError",
}

func IsTimeout(err error) bool {
	return errors.Is(err, timeoutError)
}
