package timer

import (
	"container/list"
	"context"
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
	expire   uint64
	callback func()
}

type timerNode struct {
	*list.List
}

type timer struct {
	// 单调递增累加值
	jiffies uint64

	// 256个槽位
	t1 [nearSize]timerNode
	// 4个64槽位, 代表不同的刻度
	t2Tot5 [4][levelSize]timerNode

	current int64

	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimer() *timer {
	ctx, cancel := context.WithCancel(context.Background())
	return &timer{ctx: ctx, cancel: cancel}
}

/*
func l2Max() uint64 {
	return 1 << (nearShift + levelShift)
}

func l3Max() uint64 {
	return 1 << (nearShift + 2*levelShift)
}

func l4Max() uint64 {
	return 1 << (nearShift + 3*levelShift)
}
*/

func levelMax(index int) uint64 {
	return 1 << (nearShift + index*levelShift)
}

func (t *timer) index(n int) uint64 {
	return (t.jiffies >> (nearShift + levelShift*n)) & levelMask
}

func (t *timer) add(expire time.Duration, callback func()) (node *timerNode) {
	//idx := expire - t.current
	idx := expire / (time.Millisecond * 10)
	expire = idx

	var currLevel int //debug
	var index int64

	defer func() {
		fmt.Printf("node:%p:::index:%d, currLevel:%d, idx:%d\n", node, index, currLevel, idx)
	}()

	if idx < nearSize {
		i := uint64(expire) & nearMask
		node = &t.t1[i]
		currLevel = 1
		return node
	}

	// 假如idx < 0
	// TODO

	// TODO 时间溢出

	for i := 0; i <= 3; i++ {
		if uint64(idx) < levelMax(i+1) {
			index = int64(expire) >> (nearShift + i*levelMask) & levelMask
			node = &t.t2Tot5[i][index]
			currLevel = i
			break
		}
	}

	return node
}

func (t *timer) AfterFunc(expire time.Duration, callback func()) *Time {
}

func (t *timer) Stop() {
	t.cancel()
}

func (t *timer) run() {
	t.jiffies++
}

func (t *timer) Run() {
	// 10ms精度
	tk := time.NewTimer(time.Millisecond * 10)

	for {
		select {
		case <-tk.C:
			t.run()
		case <-t.ctx.Done():
			return
		}
	}
}
