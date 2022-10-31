package timer

import "time"

type minHeapNode struct {
	stop       uint32
	callback   func()
	userExpire time.Duration
	isSchedule bool
}
