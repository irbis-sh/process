package process_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/irbis-sh/process"
)

func TestPIDExecutablePath(t *testing.T) {
	t.Parallel()

	pid := process.PID(os.Getpid())

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}

	path, err := pid.ExecutablePath()
	if err != nil {
		t.Fatalf("PID(%d).ExecutablePath(): %v", pid, err)
	}
	if path != exe {
		t.Errorf("ExecutablePath = %q, want %q", path, exe)
	}
}

func TestPIDName(t *testing.T) {
	t.Parallel()

	pid := process.PID(os.Getpid())

	name, err := pid.Name()
	if err != nil {
		t.Fatalf("PID(%d).Name(): %v", pid, err)
	}
	if name == "" {
		t.Errorf("Name() = %q, want non-empty name", name)
	}
}
