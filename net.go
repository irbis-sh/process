package process

// FindPIDBySourcePort returns the PID that owns the given TCP/IPv4 source
// port. It returns [ErrNotFound] if no process owns the port.
func FindPIDBySourcePort(port uint16) (PID, error) {
	return findPIDBySourcePort(port)
}
