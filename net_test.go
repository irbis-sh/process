package process_test

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
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

func TestFindPIDByRequest(t *testing.T) {
	t.Parallel()

	t.Run("finds owning process for local HTTP request", func(t *testing.T) {
		t.Parallel()

		type result struct {
			pid process.PID
			err error
		}

		results := make(chan result, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			pid, err := process.FindPIDByRequest(r)
			results <- result{pid: pid, err: err}
		}))
		defer srv.Close()

		resp, err := srv.Client().Get(srv.URL)
		if err != nil {
			t.Fatalf("GET %s: %v", srv.URL, err)
		}
		defer resp.Body.Close()

		got := <-results
		if got.err != nil {
			t.Fatalf("FindPIDByRequest: %v", got.err)
		}
		if int(got.pid) != os.Getpid() {
			t.Errorf("PID = %d, want %d", got.pid, os.Getpid())
		}
	})

	t.Run("returns error for malformed RemoteAddr", func(t *testing.T) {
		t.Parallel()

		_, err := process.FindPIDByRequest(&http.Request{RemoteAddr: "127.0.0.1"})
		if err == nil {
			t.Fatal("err = nil")
		}
	})

	t.Run("returns error for invalid source port", func(t *testing.T) {
		t.Parallel()

		_, err := process.FindPIDByRequest(&http.Request{RemoteAddr: "127.0.0.1:meow"})
		if err == nil {
			t.Fatal("err = nil")
		}
	})
}
