// Package process provides access to operating system process information.
//
// It defines [PID], a process identifier with methods for retrieving process
// metadata, and functions for locating processes by network port.
package process

import "errors"

// PID identifies an operating system process.
type PID int

var (
	// ErrNotFound is returned when no matching process is found.
	ErrNotFound = errors.New("process not found")
)
