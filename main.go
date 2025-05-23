package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"time"

	"github.com/mcs-unity/scheduler/internal/channel"
	"github.com/mcs-unity/scheduler/internal/fatal"
	"github.com/mcs-unity/scheduler/internal/shared"
	"github.com/mcs-unity/scheduler/pkg/schedule"
)

func test(sleep time.Duration, v string, closeCh, simPanic, cancelFn, infiniteLoop bool) schedule.Operation {
	return func(result channel.IChannel, cancel context.CancelCauseFunc) {
		defer fatal.FatalCause(cancel)

		if simPanic {
			panic("panic bro")
		}

		if cancelFn {
			cancel(errors.New("I got canceled"))
			return
		}

		if infiniteLoop {
			for {
				time.Sleep(sleep * 4)
				fmt.Println("endless leek")
			}
		}

		if closeCh {
			result.Close()
			return
		}

		time.Sleep(sleep)
		result.Send(v)
	}
}

func handle(err error) {
	log.Println(err)
}

func success(v shared.Unknown) {
	log.Println(v)
}

func main() {
	defer fatal.Capture()
	schedule.Init(5)
	if err := schedule.New("timeout", schedule.NORMAL, 2*time.Second, test(3*time.Second, "test fail timeout", false, false, false, false), handle, success); err != nil {
		fmt.Println(err)
	}

	if err := schedule.New("success", schedule.NORMAL, 1*time.Second, test(900*time.Millisecond, "test success", false, false, false, false), handle, success); err != nil {
		fmt.Println(err)
	}

	if err := schedule.New("infinite", schedule.NORMAL, 1*time.Second, test(1*time.Second, "infinite loop", false, false, false, true), handle, success); err != nil {
		fmt.Println(err)
	}

	if err := schedule.New("close channel", schedule.NORMAL, 1*time.Second, test(1*time.Second, "close channel no data", true, false, false, false), handle, success); err != nil {
		fmt.Println(err)
	}

	if err := schedule.New("panic", schedule.HIGH, 1*time.Second, test(1*time.Second, "panic no data", false, true, false, false), handle, success); err != nil {
		fmt.Println(err)
	}

	if err := schedule.New("cancel", schedule.NORMAL, 2*time.Second, test(1*time.Second, "cancel no data", false, false, true, false), handle, success); err != nil {
		fmt.Println(err)
	}

	counter := 0
	for i := range 20 {
		if counter <= 5 {
			if err := schedule.New("queue-check-high"+strconv.Itoa(i), schedule.HIGH, 2*time.Second, test(100*time.Millisecond, "test success high"+strconv.Itoa(i), false, false, false, false), handle, success); err != nil {
				fmt.Println(err)
			}
			counter++
		} else {
			counter = 0
		}

		if i == 20/2 {
			if err := schedule.New("queue-check"+strconv.Itoa(i), schedule.REALTIME, 2*time.Second, test(50*time.Millisecond, "test success realtime"+strconv.Itoa(i), false, false, false, false), handle, success); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := schedule.New("queue-check"+strconv.Itoa(i), schedule.NORMAL, 2*time.Second, test(1000*time.Millisecond, "test success normal"+strconv.Itoa(i), false, false, false, false), handle, success); err != nil {
				fmt.Println(err)
			}
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	fmt.Println(<-ch)
}

func init() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile)
}
