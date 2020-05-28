package timer

import (
	"container/list"
	"time"
)

const (
	nearShift  = 8
	nearSize   = 1 << nearShift
	levelShift = 6
	levelSize  = 1 << levelShift
	nearMask   = nearSize - 1
	levelMask  = levelSize - 1
)

type timerSession struct {
	expire   int64 //单位精确到ms
	callback func()
}

type timerNode struct {
	*list.List
}

type timer struct {
	// 256个槽位
	t1 [nearSize]timerNode
	//4个64槽位, 代表不同的刻度
	t2 [levelSize]timerNode
	t3 [levelSize]timerNode
	t4 [levelSize]timerNode
	t5 [levelSize]timerNode

	current int64
}

func NewTimer() *timer {
	return &timer{}
}

func (t *timer) add(expire int64, callback func()) {
	idx := expire - t.current

	var node timerNode
	if idx < nearSize {
		node = t.t1[expire&nearMask]
	} else if idx < 1<<nearShift+levelShift {
		i := expire >> nearShift & levelMask
		node = t.t2[i]
	} else if idx < 1<<nearShift+2*levelShift {
		i := expire >> (nearShift + levelShift) & levelMask
		node = t.t3[i]
	} else if idx < 1<<nearShift+3*levelShift {
		i := expire >> (nearShift + 2*levelShift) & levelMask
		node = t.t4[i]
	} else if idx < 0 {
		//TODO
	} else {
		i := expire >> (nearShift + 3*levelShift) & levelMask
		node = t.t5[i]
	}

	_ = node
}

func (t *timer) Add(expire int64, callback func()) {
}
