package bench

import (
	"bytes"
	"io"
	"os"
)

// withStdoutCapture runs fn with os.Stdout redirected; returns captured bytes and fn's error.
// fn must run on the calling goroutine (required for raylib main-thread rules).
func withStdoutCapture(fn func() error) ([]byte, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = w
	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(&buf, r)
		r.Close()
		close(done)
	}()
	execErr := fn()
	w.Close()
	<-done
	os.Stdout = old
	return buf.Bytes(), execErr
}
