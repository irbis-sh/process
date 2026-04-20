package process

// Name returns a best-effort display name for the process.
//
// Its meaning is operating system dependent. The returned name may come from
// process metadata rather than the executable filename, may be truncated, and
// is not guaranteed to be stable over the lifetime of the process.
//
// Name is intended for display and diagnostics. Callers that need a stable
// identifier should use PID or ExecutablePath instead.
func (pid PID) Name() (string, error) {
	return pidName(pid)
}

// ExecutablePath returns the filesystem path to the process executable.
func (pid PID) ExecutablePath() (string, error) {
	return pidExecutablePath(pid)
}
