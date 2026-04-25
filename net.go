package process

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

// FindPIDBySourcePort returns the PID that owns the given TCP/IPv4 source
// port. It returns [ErrNotFound] if no process owns the port.
func FindPIDBySourcePort(port uint16) (PID, error) {
	return findPIDBySourcePort(port)
}

// FindPIDByRequest returns the PID that owns the TCP/IPv4 source port for r.
//
// Only works for local requests. Returns [ErrNotFound] if no process owns the port.
func FindPIDByRequest(r *http.Request) (PID, error) {
	_, sourcePort, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return 0, fmt.Errorf("parse RemoteAddr: %v", err)
	}
	sourcePortNum, err := strconv.ParseUint(sourcePort, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("parse source port: %v", err)
	}

	return FindPIDBySourcePort(uint16(sourcePortNum))
}
