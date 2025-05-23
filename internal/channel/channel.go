package channel

import (
	"errors"
	"sync"

	"github.com/mcs-unity/scheduler/internal/shared"
)

func (c *Channel) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return
	}

	close(c.channel)
	c.close = true
}

func (c *Channel) Send(data shared.Unknown) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return errors.New(ErrClosed)
	}

	c.channel <- data
	return nil
}

func (c Channel) Wait() <-chan shared.Unknown {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.close {
		return nil
	}
	return c.channel
}

func New() IChannel {
	return &Channel{
		lock:    &sync.Mutex{},
		channel: make(chan shared.Unknown, 1),
		close:   false,
	}
}
