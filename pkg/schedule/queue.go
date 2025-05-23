package schedule

import (
	"slices"
)

//! Issues
/*
Make a test and ensure that the queue will be switched between high and normal
*/

func (sc scheduler) priority() Priority {
	if sc.listCounter < maxNormalWait && len(sc.high) > 0 {
		return HIGH
	}

	return NORMAL
}

func (sc *scheduler) getList(p Priority) []Pending {
	if p == HIGH {
		return sc.high
	}

	sc.listCounter = 0
	return sc.normal
}

func (sc *scheduler) removeItem(p Priority, index int) {
	if p == HIGH {
		sc.high = slices.Delete(sc.high, index, index+1)
	} else {
		sc.normal = slices.Delete(sc.normal, index, index+1)
	}
}

func (sc *scheduler) next() {
	p := sc.priority()
	list := sc.getList(p)
	var next *Pending
	var index = -1
	if len(list) == 0 {
		return
	}

	for i, v := range list {
		if next != nil && next.p == HIGH {
			break
		}
		next = &v
		index = i
	}

	if next == nil {
		return
	}

	sc.removeItem(p, index)
	if p == HIGH {
		sc.listCounter++
	}

	sc.list <- next.s
}
