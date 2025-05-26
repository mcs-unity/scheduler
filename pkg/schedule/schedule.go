package schedule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/mcs-unity/scheduler/internal/shared"

	"time"

	"github.com/mcs-unity/scheduler/internal/channel"
	"github.com/mcs-unity/scheduler/internal/fatal"
)

var monitor *scheduler

func (s *scheduler) queue(p Priority, sc Schedule) {
	op := Pending{p, sc}
	if p == HIGH {
		s.high = append(s.high, op)
	} else {
		s.normal = append(s.normal, op)
	}
}

func New(i ID, p Priority, t Timeout, f Operation, e ErrorCallback, s SuccessCallback) error {
	monitor.lock.Lock()
	defer monitor.lock.Unlock()

	if monitor.closed {
		return errors.New(ErrSchedulerClosing)
	}

	if err := validate(i, t, f, e); err != nil {
		return err
	}

	sc := Schedule{
		I: i,
		T: t,
		F: f,
		S: s,
		E: e,
	}

	if monitor.len() > monitor.concurrency && p != REALTIME {
		monitor.queue(p, sc)
	} else {
		monitor.add()
		monitor.list <- sc
	}
	return nil
}

func throwError(id string, reason string) error {
	return fmt.Errorf(errFormat, id, reason)
}

func handle(sc Schedule, v shared.Unknown, start time.Time, ok bool) {
	if !ok {
		sc.E(throwError(sc.I, ErrChannelClosed))
		return
	}

	endtime(sc.I, start, sc.T)
	if sc.S != nil {
		sc.S(v)
	}
}

func execute(sc Schedule) {
	defer fatal.FatalCapture(sc.I)
	defer monitor.done()

	ctx, cancel := context.WithTimeout(context.Background(), sc.T)
	ctxCause, cancelCause := context.WithCancelCause(ctx)
	defer cancel()

	start := time.Now()
	result := channel.New()
	defer result.Close()
	go sc.F(result, cancelCause)

	select {
	case v, ok := <-result.Wait():
		handle(sc, v, start, ok)
	case <-ctxCause.Done():
		sc.E(throwError(sc.I, context.Cause(ctxCause).Error()))
	}
}

func start() {
	log.Println("starting monitor")
	defer fatal.Capture()
	defer close(monitor.list)
	defer close(monitor.kill)
loop:
	for {
		select {
		case sc, ok := <-monitor.list:
			if !ok {
				break loop
			}
			go execute(sc)
		case <-monitor.kill:
			break loop
		}
	}
}

func Kill() {
	monitor.closed = true
	monitor.lock.Lock()
	defer monitor.lock.Unlock()
	monitor.kill <- struct{}{}
}

func Init(c uint64) {
	if monitor != nil && !monitor.closed {
		return
	}

	monitor = &scheduler{
		lock:        &sync.Mutex{},
		list:        make(chan Schedule, 1),
		kill:        make(chan struct{}, 1),
		closed:      false,
		concurrency: int(c),
		counter:     0,
		normal:      make([]Pending, 0),
		high:        make([]Pending, 0),
	}
	go start()
}
