package schedule

import (
	"errors"
	"strings"
	"time"
)

func validate(i ID, t Timeout, f Operation, e ErrorCallback) error {
	if strings.Trim(i, "") == "" {
		return errors.New(ErrIdIsEmpty)
	}

	if t > 5*time.Second {
		return errors.New(ErrTimeoutIsGreater)
	}

	if f == nil {
		return errors.New(ErrOperationIsNil)
	}

	if e == nil {
		return errors.New(ErrCallbackIsNil)
	}

	return nil
}
