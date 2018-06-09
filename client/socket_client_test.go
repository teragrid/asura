package asuracli_test

import (
	"errors"
	"testing"
	"time"

	"github.com/teragrid/asura/client"
)

func TestSocketClientStopForErrorDeadlock(t *testing.T) {
	c := asuracli.NewSocketClient(":80", false)
	err := errors.New("foo-teragrid")

	// See Issue https://github.com/teragrid/asura/issues/114
	doneChan := make(chan bool)
	go func() {
		defer close(doneChan)
		c.StopForError(err)
		c.StopForError(err)
	}()

	select {
	case <-doneChan:
	case <-time.After(time.Second * 4):
		t.Fatalf("Test took too long, potential deadlock still exists")
	}
}
