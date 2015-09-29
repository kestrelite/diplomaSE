// Package util provides useful functions and objects for developing diplomaSE.
package util

import (
	"io/ioutil"
	"log"
	"os"
)

type tracer struct {
	*log.Logger
	activated bool
}

// The GetNewTracer function returns a predefined logger object with minimal
// formatting. Messages to the new logger are discarded by default until and
// unless it is explicitly "turned on" using Activate() or Toggle().
func GetNewTracer() *tracer {
	return &tracer{log.New(ioutil.Discard, "", 0), false}
}

// The Activate method tells the tracer to route messages to stdout.
func (t *tracer) Activate() {
	t.SetOutput(os.Stdout)
	t.activated = true
}

/// The Deactivate method tells the tracer to discard all messages it receives.
func (t *tracer) Deactivate() {
	t.SetOutput(ioutil.Discard)
	t.activated = false
}
