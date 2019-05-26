package circuit

import (
	"context"
	"errors"
	"time"
)

type State int

const (
	Unknown State = iota
	Failure
	Success
)

type Counter interface {
	Count(State)
	ConsecutiveFailures() uint32
	LastActivity() time.Time
	Reset()
}

type CircuitCounter int

func (c CircuitCounter) Count(state State) {

}

func (c CircuitCounter) ConsecutiveFailures() uint32 {
	return 1
}

func (c CircuitCounter) LastActivity() time.Time {
	return time.Now()
}

func (c CircuitCounter) Reset() {

}

func NewCounter() Counter {
	var c CircuitCounter = 1
	return c
}

type Circuit func(context.Context) error

var (
	ErrServiceUnavailable = errors.New("service unavailable")
)

func Breaker(c Circuit, failureThreshold uint32) Circuit {
	cnt := NewCounter()

	return func(ctx context.Context) error {
		if cnt.ConsecutiveFailures() >= failureThreshold {
			canRetry := func(cnt Counter) bool {
				backoffLevel := cnt.ConsecutiveFailures() - failureThreshold

				// Calculates when should the circuit breaker resume propagating requests
				// to the service
				shouldRetryAt := cnt.LastActivity().Add(time.Second * 2 << backoffLevel)

				return time.Now().After(shouldRetryAt)
			}

			if !canRetry(cnt) {
				return ErrServiceUnavailable
			}
		}

		if err := c(ctx); err != nil {
			cnt.Count(Failure)
			return err
		}

		cnt.Count(Success)
		return nil
	}
}
