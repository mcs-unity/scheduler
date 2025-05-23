package channel

import (
	"sync"

	"github.com/mcs-unity/scheduler/internal/shared"
)

const (
	ErrClosed = "Channel closed"
)

type IChannel interface {
	Close()
	Send(data shared.Unknown) error
	Wait() <-chan shared.Unknown
}

type Channel struct {
	lock    sync.Locker
	channel chan shared.Unknown
	close   bool
}
