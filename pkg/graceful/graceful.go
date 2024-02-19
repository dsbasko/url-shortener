package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

// Graceful is a struct that encapsulates a WaitGroup.
// It is used to wait for all added goroutines to finish.
type Graceful struct {
	wg        sync.WaitGroup
	total     atomic.Int32
	completed atomic.Int32
}

// Global variables for the Graceful instance and a Once object to ensure it's only initialized once.
var (
	once     sync.Once
	graceful *Graceful
)

// init is a special function that gets called when the package is imported.
// It initializes the graceful object once.
func init() { //nolint:gochecknoinits
	once.Do(func() {
		graceful = &Graceful{}
	})
}

// Context creates a new context that is canceled when any of the provided signals are received.
// It wraps the signal.NotifyContext function, taking a parent context and a variable number of os.Signal values.
// It returns the new context and its associated cancel function.
func Context(parentCtx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parentCtx, signals...)
}

// Add increments the WaitGroup counter by one.
// This should be called before starting a goroutine.
func Add() {
	graceful.total.Add(1)
	graceful.wg.Add(1)
}

// Done decrements the WaitGroup counter by one.
// This should be called when a goroutine finishes.
func Done() {
	graceful.completed.Add(1)
	graceful.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero.
// This should be called where you want to wait for all goroutines to finish.
func Wait() {
	graceful.wg.Wait()
}

// Count returns a string representing the current state of the Graceful object.
// The format of the returned string is: "completed/total".
// This can be useful for tracking the progress of goroutine execution.
func Count() string {
	return fmt.Sprintf("%v/%v", graceful.completed.Load(), graceful.total.Load())
}

// Reset reinitializes the state of the Graceful object.
// It resets the WaitGroup and sets the total and completed counters to zero.
// This function can be useful when you want to reuse the Graceful object for a new set of goroutines.
func Reset() {
	graceful.wg = sync.WaitGroup{}
	graceful.total.Store(0)
	graceful.completed.Store(0)
}

// CleanFn executes a cleanup function after a specified duration.
// It uses a goroutine to run the cleanup function and a timer to control the duration.
// Once the cleanup function is done, it signals the doneCh channel.
// If the timer expires before the cleanup function finishes, the function returns immediately.
func CleanFn(duration time.Duration, cleanFn func()) {
	defer Done()
	doneCh := make(chan struct{})

	t := time.NewTimer(duration)

	go func() {
		cleanFn()
		close(doneCh)
	}()

	select {
	case <-t.C:
		return
	case <-doneCh:
		return
	}
}
