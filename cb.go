package cb

import (
	"errors"
	"time"
)

type State int

const (
	StateClose State = iota
	StateOpen
	StateHalfOpen
)

type Breaker struct {
	failureThreshold int
	resetAfter       time.Duration
	lastFailedAt     time.Time
	failures         int
	Timeout          int
}

// State
func (b *Breaker) State() State {
	switch {
	case b.failures > b.failureThreshold && time.Since(b.lastFailedAt) > b.resetAfter:
		return StateHalfOpen
	case b.failures > b.failureThreshold:
		return StateOpen
	}

	return StateClose
}

func (b *Breaker) Call(fn func() error) error {
	if b.State() == StateClose || b.State() == StateHalfOpen {
		err := fn()
		if err != nil {
			b.failures += 1
			b.lastFailedAt = time.Now()
		} else {
			b.failures = 0
		}

		return err
	}

	return errors.New("")
}
