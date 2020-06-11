package timer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

const (
	haveStop = uint32(1)
	//stopGrab = 1 << (iota + 1)
	//pushGrab
)

// 先使用sync.Mutex实现功能
// 后面使用cas优化
type Time struct {
	*list.List

	sync.Mutex
}

func (t *Time) lockPushBack(node *timeNode) {
	t.Lock()
	defer t.Unlock()
	if atomic.LoadUint32(&node.lock) == haveStop {
		return
	}

	node.element = t.PushBack(node)
	node.list = t
}

type timeNode struct {
	expire     uint64
	userExpire time.Duration
	callback   func()
	isSchedule bool
	close      uint32
	lock       uint32

	list    *Time
	element *list.Element
}

/*
func (t *timeNode) grab() {
	for {
		prevVal := atomic.LoadUint32(&t.lock)
		if atomic.CompareAndSwapUint32(&t.lock, prevVal, stopGrab) {
			break
		}
	}
}
*/

func (t *timeNode) Stop() {
	//这里和32行是竞争关系，拷贝一个副本，防止出现unlock unlock的情况
	cpyList := t.list
	cpyList.Lock()
	defer cpyList.Unlock()

	atomic.StoreUint32(&t.close, haveStop)

	//t.grab()

	t.list.Remove(t.element)
}
