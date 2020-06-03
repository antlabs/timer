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

type timer struct {
	// 单调递增累加值, 走过一个时间片就+1
	jiffies uint64

	// 256个槽位
	t1 [nearSize]*Time
	// 4个64槽位, 代表不同的刻度
	t2Tot5 [4][levelSize]*Time

	// 时间只精确到10ms
	// curTimePoint 为1就是10ms 为2就是20ms
	curTimePoint time.Duration
	// 上下文
	ctx context.Context
	// 取消函数
	cancel context.CancelFunc
}

func NewTimer() *timer {
	ctx, cancel := context.WithCancel(context.Background())
	t := &timer{ctx: ctx, cancel: cancel}
	t.init()
	return t
}

func (t *timer) init() {
	for i := 0; i < nearSize; i++ {
		t.t1[i] = &Time{List: list.New()}
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < levelSize; j++ {
			t.t2Tot5[i][j] = &Time{List: list.New()}
		}
	}

	t.curTimePoint = get10Ms()
}

func levelMax(index int) uint64 {
	return 1 << (nearShift + index*levelShift)
}

func (t *timer) index(n int) uint64 {
	return (t.jiffies >> (nearShift + levelShift*n)) & levelMask
}

func (t *timer) add(expire time.Duration, callback func()) (node *Time) {
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
		node = t.t1[i]
		currLevel = 1
		return node
	}

	// 假如idx < 0
	// TODO

	// TODO 时间溢出

	for i := 0; i <= 3; i++ {
		if uint64(idx) < levelMax(i+1) {
			index = int64(expire) >> (nearShift + i*levelMask) & levelMask
			node = t.t2Tot5[i][index]
			currLevel = i
			break
		}
	}

	return node
}

func (t *timer) AfterFunc(expire time.Duration, callback func()) *Time {
	return t.add(expire, callback)
}

func (t *timer) Stop() {
	t.cancel()
}

func (t *timer) moveAndExec() {
	//1. 先移动到near链表里面
	//2. 再执行
	t.jiffies++
}

func (t *timer) run() {
	// 先判断是否需要更新
	// 内核里面实现使用了全局jiffies和本地的jiffies比较,应用层没有jiffies，直接使用时间比较
	// 这也是skynet里面的做法

	ms10 := get10Ms()

	if ms10 < t.curTimePoint {

		fmt.Printf("github.com/antlabs/timer:Time has been called back?from(%d)(%d)\n",
			ms10, t.curTimePoint)

		t.curTimePoint = ms10
		return
	}

	diff := ms10 - t.curTimePoint
	t.curTimePoint = ms10
	for i := 0; i < int(diff); i++ {
		t.moveAndExec()
	}

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
