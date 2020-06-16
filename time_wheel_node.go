package timer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/antlabs/stl/list"
)

const (
	haveStop = uint32(1)
	//stopGrab = 1 << (iota + 1)
	//pushGrab
)

// 先使用sync.Mutex实现功能
// 后面使用cas优化
type Time struct {
	timeNode
	sync.Mutex
}

func newTimeHead() *Time {
	head := &Time{}
	head.Init()
	return head
}

func (t *Time) lockPushBack(node *timeNode) {
	t.Lock()
	defer t.Unlock()
	if atomic.LoadUint32(&node.lock) == haveStop {
		return
	}

	t.AddTail(&node.Head)
	node.list = t
}

type timeNode struct {
	expire     uint64
	userExpire time.Duration
	callback   func()
	isSchedule bool
	close      uint32
	lock       uint32

	list *Time

	list.Head
}

func (t *timeNode) Stop() {
	//这里为什么修改成cpyList := t.list
	//如果直接使用t.list.Lock()和t.list.Unlock()就会和lockPushBack函数里的node.list = t 是竞争关系
	//lockPushBack拿到锁，修改国t.list的值。这时候Stop函数里面的t.list.Lock()持有旧链表里的锁。t.list.unlock新链表里的锁，发触发unlock unlock情况。
	
	//TODO：思考有没有新的竞争关系。。。
	cpyList := t.list
	cpyList.Lock()
	defer cpyList.Unlock()

	atomic.StoreUint32(&t.close, haveStop)

	t.list.Del(&t.Head)
}
