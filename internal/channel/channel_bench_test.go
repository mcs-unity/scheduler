package channel

import (
	"context"
	"testing"
	"time"
)

// Benchmark test will create a new channel in every loop
// the goal is to see how this approach will impact an environment
// where a new channel will need to be created with every task.
func BenchmarkChannel(b *testing.B) {
	val := "test"
	fn := func(channel IChannel) {
		if err := channel.Send(val); err != nil {
			b.Error(err)
		}
	}

	for b.Loop() {
		channel := New()
		defer channel.Close()
		ch := channel.Wait()
		go fn(channel)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			b.Error("timeout")
		case r := <-ch:
			if r != val {
				b.Errorf("%s did not match expected %s", r, val)
			}
		}
	}
}

// here we try to make the channel a bit more usable by trying
// to reuse the channel and see how it behaves in terms of
// resource consumption, this approach is best should you be
// restrained by a low amount of recourses.
func BenchmarkChannelReusability(b *testing.B) {
	channel := New()
	defer channel.Close()
	ch := channel.Wait()
	val := "test"
	fn := func() {
		if err := channel.Send(val); err != nil {
			b.Error(err)
		}
	}

	for b.Loop() {
		go fn()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			b.Error("timeout")
		case r := <-ch:
			if r != val {
				b.Errorf("%s did not match expected %s", r, val)
			}
		}
	}
}
