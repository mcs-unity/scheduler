package fatal

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

const expected = "got %s expected %s"

func TestFatalCausePanic(t *testing.T) {
	initial := "failed"
	ctx, cancel := context.WithTimeoutCause(context.Background(), 300*time.Millisecond, errors.New(initial))
	defer cancel()
	ctx, cancelCause := context.WithCancelCause(ctx)
	reason := "canceled"
	go func() {
		defer FatalCause(cancelCause)
		panic(reason)
	}()

	<-ctx.Done()
	cancelReason := context.Cause(ctx).Error()
	if !strings.Contains(cancelReason, reason) {
		t.Errorf(expected, cancelReason, reason)
	}
}

func TestHandlePanic(t *testing.T) {
	errorMessage := "handle me"
	message := handlePanic(errorMessage)

	if !strings.Contains(message, errorMessage) {
		t.Errorf(expected, message, errorMessage)
	}

	err := errors.New("i'm an error")
	message = handlePanic(err)

	if !strings.Contains(message, err.Error()) {
		t.Errorf(expected, message, err.Error())
	}

	message = handlePanic([]string{})

	if !strings.Contains(message, unknownErr) {
		t.Errorf(expected, message, unknownErr)
	}
}
