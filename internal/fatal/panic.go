package fatal

import (
	"context"
	"fmt"

	"github.com/mcs-unity/scheduler/internal/shared"
)

func handlePanic(err shared.Unknown) string {
	message := ""
	switch e := err.(type) {
	case string:
		message = e
	case error:
		message = e.Error()
	default:
		message = unknownErr
	}

	return fmt.Sprintf(format, message)
}

func Capture() {
	if capture := recover(); capture != nil {
		fmt.Printf(formatWithReason+"\n", handlePanic(capture))
	}
}

func FatalCapture(id string) {
	if capture := recover(); capture != nil {
		fmt.Printf(formatWithId, id, handlePanic(capture))
	}
}

func FatalCause(cancel context.CancelCauseFunc) {
	if capture := recover(); capture != nil {
		cancel(fmt.Errorf(formatWithReason, handlePanic(capture)))
	}
}
