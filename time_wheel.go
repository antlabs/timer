package timer

import (
	"container/list"
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	nearShift = 8

	nearSize = 1 << nearShift

	levelShift = 6

	levelSize = 1 << levelShift

	nearMask = nearSize - 1

	levelMask = levelSize - 1
)

type timeWheel struct {
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

func newTimeWheel() *timeWheel {

	ctx, cancel := context.WithCancel(context.Background())

	t := &timeWheel{ctx: ctx, cancel: cancel}

	t.init()

	return t
}

func (t *timeWheel) init() {

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

func maxVal() uint64 {
	return (1 << (nearShift + 4*levelShift)) - 1
}

func levelMax(index int) uint64 {
	return 1 << (nearShift + index*levelShift)
}

func (t *timeWheel) index(n int) uint64 {
	return (t.jiffies >> (nearShift + levelShift*n)) & levelMask
}

func (t *timeWheel) add(node *timeNode, jiffies uint64) *timeNode {

	var head *Time
	expire := node.expire
	idx := expire - jiffies

	if idx < nearSize {

		i := uint64(expire) & nearMask
		head = t.t1[i]

	} else {

		max := maxVal()
		for i := 0; i <= 3; i++ {

			if idx > max {
				idx = max
				expire = idx + jiffies
			}

			if uint64(idx) < levelMax(i+1) {
				index := int64(expire) >> (nearShift + i*levelMask) & levelMask
				head = t.t2Tot5[i][index]
				break
			}
		}
		//fmt.Printf("idx:%d i:%p\n", idx, head)
	}

	if head == nil {
		panic("not found head")
	}

	node.element = head.PushBack(node)

	return node
}

func (t *timeWheel) AfterFunc(expire time.Duration, callback func()) TimeNoder {

	jiffies := atomic.LoadUint64(&t.jiffies)

	expire = expire/(time.Millisecond*10) + time.Duration(jiffies)

	node := &timeNode{
		expire:   uint64(expire),
		callback: callback,
	}

	return t.add(node, jiffies)
}

func (t *timeWheel) getExpire(expire time.Duration, jiffies uint64) time.Duration {
	return expire/(time.Millisecond*10) + time.Duration(jiffies)
}

func (t *timeWheel) ScheduleFunc(userExpire time.Duration, callback func()) TimeNoder {

	jiffies := atomic.LoadUint64(&t.jiffies)

	expire := t.getExpire(userExpire, jiffies)

	node := &timeNode{
		userExpire: userExpire,
		expire:     uint64(expire),
		callback:   callback,
		isSchedule: true,
	}

	return t.add(node, jiffies)
}

func (t *timeWheel) Stop() {
	t.cancel()
}

// 移动链表
func (t *timeWheel) cascade(levelIndex int, index int) {

	tmp := list.New()

	l := t.t2Tot5[levelIndex][index]

	tmp.PushBackList(l.List)

	t.t2Tot5[levelIndex][index].Init()

	for e := tmp.Front(); e != nil; e = e.Next() {
		t.add(e.Value.(*timeNode), atomic.LoadUint64(&t.jiffies))
	}
}

// moveAndExec函数功能
//1. 先移动到near链表里面
//2. near链表节点为空时，从上一层里面移动一些节点到下一层
//3. 再执行
func (t *timeWheel) moveAndExec() {

	// 这里时间溢出
	if uint32(t.jiffies) == 0 {
		// TODO
		// return
	}

	//如果本层的盘子没有定时器，这时候和上层的盘子移动一些过来
	index := t.jiffies & nearMask
	if index == 0 {
		for i := 0; i <= 3; i++ {
			index = t.index(i)
			if index != 0 {
				t.cascade(i, int(index))
				break
			}
		}
	}

	// t.jiffies 只在run go程里面修改，所以不需要加atomic函数原子取
	t.jiffies++

	if t.t1[index].List.Len() == 0 {
		return
	}

	// 执行
	head := Time{List: list.New()}

	head.PushBackList(t.t1[index].List)
	t.t1[index].List.Init()

	for e := head.Front(); e != nil; {

		val := e.Value.(*timeNode)
		next := e.Next()
		head.Remove(e)
		e = next

		go val.callback()

		if val.isSchedule {
			jiffies := t.jiffies
			val.expire = uint64(t.getExpire(val.userExpire, jiffies))
			t.add(val, jiffies)
		}
	}
}

func (t *timeWheel) run() {
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

func (t *timeWheel) Run() {

	// 10ms精度
	tk := time.NewTicker(time.Millisecond * 10)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			t.run()
		case <-t.ctx.Done():
			return
		}
	}
}
