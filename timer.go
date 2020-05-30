package timer

import (
	"container/list"
	"fmt"
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

func (t *timer) debug() {
	fmt.Printf("t1 address:%p\n", &t.t1)
	fmt.Printf("t2 address:%p\n", &t.t2)
	fmt.Printf("t3 address:%p\n", &t.t3)
	fmt.Printf("t4 address:%p\n", &t.t4)
	fmt.Printf("t5 address:%p\n", &t.t5)
}

func l2Max() int64 {
	return 1 << (nearShift + levelShift)
}

func l3Max() int64 {
	return 1 << (nearShift + 2*levelShift)
}

func l4Max() int64 {
	return 1 << (nearShift + 3*levelShift)
}

func (t *timer) add(expire time.Duration, callback func()) *timerNode {
	//idx := expire - t.current
	idx := expire / (time.Millisecond * 10)
	expire = idx

	var node *timerNode
	var i int64       //debug
	var currLevel int //debug

	if idx < nearSize {
		i = int64(expire) & nearMask
		node = &t.t1[i]
		currLevel = 1
	} else if int64(idx) < l2Max() {
		i = int64(expire) >> nearShift & levelMask
		node = &t.t2[i]
		currLevel = 2
	} else if int64(idx) < l3Max() {
		i = int64(expire) >> (nearShift + levelShift) & levelMask
		node = &t.t3[i]
		currLevel = 3
	} else if int64(idx) < l4Max() {
		i = int64(expire) >> (nearShift + 2*levelShift) & levelMask
		node = &t.t4[i]
		currLevel = 4
	} else if idx < 0 {
		//TODO
	} else {
		i = int64(expire) >> (nearShift + 3*levelShift) & levelMask
		node = &t.t5[i]
		currLevel = 5
	}

	fmt.Printf("node:%p:::index:%d, currLevel:%d, idx:%d\n", node, i, currLevel, idx)
	return node
}

func (t *timer) Add(expire int64, callback func()) {
}
