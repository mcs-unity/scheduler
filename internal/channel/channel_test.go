package channel

import (
	"context"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	channel := New()
	defer channel.Close()
	ch := channel.Wait()

	val := "test"
	go func() {
		if err := channel.Send(val); err != nil {
			t.Error(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("timeout")
	case r := <-ch:
		if r != val {
			t.Errorf("%s did not match expected %s", r, val)
		}
	}
}

func TestClosedChannel(t *testing.T) {
	channel := New()
	ch := channel.Wait()
	go func() {
		time.Sleep(100 * time.Millisecond)
		channel.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("failed to close channel")
		}
	case <-ctx.Done():
		t.Error("timeout")
	}
}

func TestCloseTwice(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	channel := New()
	channel.Close()
	channel.Close()
}

func TestSendOnClosedChannel(t *testing.T) {
	channel := New()
	channel.Close()
	if err := channel.Send("test"); err == nil {
		t.Error("should return error")
	}
}

func TestWaitIsNilWhenClosed(t *testing.T) {
	channel := New()
	channel.Close()
	if w := channel.Wait(); w != nil {
		t.Error("channel is closed wait should return nil")
	}
}
