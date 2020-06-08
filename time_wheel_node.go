package timer

import (
	"container/list"
	"sync"
	"time"
)

// 先使用sync.Mutex实现功能
// 后面使用cas优化
type Time struct {
	*list.List

	sync.Mutex
}

type timeNode struct {
	expire     uint64
	userExpire time.Duration
	callback   func()
	isSchedule bool

	list    *Time
	element *list.Element
}

func (t *timeNode) Stop() {
	t.list.Lock()
	defer t.list.Unlock()
	t.list.Remove(t.element)
}
