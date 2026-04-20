package process_test

import (
	"errors"
	"net"
	"os"
	"testing"

	"github.com/irbis-sh/process"
)

func TestFindPIDBySourcePort(t *testing.T) {
	t.Parallel()

	t.Run("finds owning process for active connection", func(t *testing.T) {
		t.Parallel()

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		serverConn := make(chan net.Conn, 1)
		go func() {
			c, _ := ln.Accept()
			serverConn <- c
		}()

		conn, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		if sc := <-serverConn; sc != nil {
			defer sc.Close()
		}

		port := conn.LocalAddr().(*net.TCPAddr).Port

		pid, err := process.FindPIDBySourcePort(uint16(port)) // #nosec G115 -- port will fit in uint16
		if err != nil {
			t.Fatalf("FindPIDBySourcePort(%d): %v", port, err)
		}

		if int(pid) != os.Getpid() {
			t.Errorf("PID = %d, want %d", pid, os.Getpid())
		}
	})

	t.Run("returns ErrNotFound for unbound port", func(t *testing.T) {
		t.Parallel()

		if _, err := process.FindPIDBySourcePort(0); !errors.Is(err, process.ErrNotFound) {
			t.Errorf("err = %v, want %v", err, process.ErrNotFound)
		}
	})
}
