package gsync

import (
	"sync"
)

// Shutdown is some generic boilerplate code to handle cleanup/stopping of some
// process that requires such handling.
// It has two signals:
//  1. Shutdown: a signal to indicate that the process should stop.
//  2. ShutdownComplete: a signal to indicate that all cleanup process has run.
type Shutdown struct {
	shutdownChan         chan struct{}
	shutdownOnce         sync.Once
	shutdownCompleteChan chan struct{}
	shutdownCompleteOnce sync.Once
	cleanups             []func()
}

// NewShutdown creates a new Shutdown instance.
func NewShutdown() Shutdown {
	return Shutdown{
		shutdownChan:         make(chan struct{}),
		shutdownCompleteChan: make(chan struct{}),
	}
}

// AddCleanup adds a cleanup function to be called when the shutdown signal is received.
func (s *Shutdown) AddCleanup(fn func()) {
	s.cleanups = append(s.cleanups, fn)
}

// Shutdown sends the shutdown signal and runs all cleanup functions.
func (s *Shutdown) Shutdown() {
	s.shutdownOnce.Do(func() {
		for _, fn := range s.cleanups {
			fn()
		}
		close(s.shutdownChan)
	})
}

// ShutdownSignal returns a channel that will be closed when the shutdown signal is received.
func (s *Shutdown) ShutdownSignal() <-chan struct{} {
	return s.shutdownChan
}

// ShutdownCompleteSignal returns a channel that will be closed when all cleanup functions have been run.
func (s *Shutdown) ShutdownCompleteSignal() <-chan struct{} {
	return s.shutdownCompleteChan
}

// SignalStopComplete signals that the shutdown process is complete and closes the shutdownCompleteChan.
func (s *Shutdown) SignalStopComplete() {
	s.shutdownCompleteOnce.Do(func() {
		close(s.shutdownCompleteChan)
	})
}
