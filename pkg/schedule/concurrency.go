package schedule

import (
	"log"
	"time"
)

func endtime(id string, start time.Time, timeout time.Duration) {
	executionTime := time.Since(start)
	if executionTime > timeout/10*8 {
		log.Printf(executionTimeFormat, id, executionTime, timeout)
	}
}

func (s *scheduler) add() {
	s.counter++
}

func (s *scheduler) done() {
	s.counter = s.counter - 1
	s.next()
}

func (s scheduler) len() int {
	return s.counter
}
