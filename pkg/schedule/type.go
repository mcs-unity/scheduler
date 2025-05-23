package schedule

import (
	"context"
	"sync"

	"time"

	"github.com/mcs-unity/scheduler/internal/channel"
	"github.com/mcs-unity/scheduler/internal/shared"
)

const max uint8 = 255
const maxErr uint8 = 2
const maxNormalWait uint8 = 1

const errFormat = "id:%s reason: %s"
const executionTimeFormat = "warning execution time id: %s time: %s timeout %s"

type ID = string
type Timeout = time.Duration
type Priority uint8

type Operation func(result channel.IChannel, cancel context.CancelCauseFunc)
type SuccessCallback func(d shared.Unknown)
type ErrorCallback func(err error)

const (
	ErrIdIsEmpty        = "id can't be an empty string"
	ErrTimeoutIsGreater = "timeout can't be greater than 5 seconds"
	ErrOperationIsNil   = "operation can't be a nil pointer"
	ErrCallbackIsNil    = "errorCallback can't be a nil pointer"
	ErrSchedulerClosing = "scheduler is closing down"
	ErrChannelClosed    = "channel closed"
)

const (
	REALTIME Priority = iota
	HIGH     Priority = iota
	NORMAL   Priority = iota
)

type Schedule struct {
	I ID
	T Timeout
	F Operation
	S SuccessCallback
	E ErrorCallback
}

type Pending struct {
	p Priority
	s Schedule
}

type scheduler struct {
	list        chan Schedule
	kill        chan struct{}
	closed      bool
	lock        sync.Locker
	concurrency int
	counter     int
	normal      []Pending
	high        []Pending
	listCounter uint8
}
