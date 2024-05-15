// Copyright 2020-2024 guonaihong, antlabs. All rights reserved.
//
// mit license
package timer

import (
	"container/heap"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

var _ Timer = (*minHeap)(nil)

var defaultTimeout = time.Hour

type minHeap struct {
	mu sync.Mutex
	minHeaps
	chAdd    chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
	wait     sync.WaitGroup
	tm       *time.Timer
	runCount int32 // 单元测试时使用
}

// 一次性定时器
func (m *minHeap) AfterFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, nil, callback, false)
}

// 周期性定时器
func (m *minHeap) ScheduleFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, nil, callback, true)
}

// 自定义下次的时间
func (m *minHeap) CustomFunc(n Next, callback func()) TimeNoder {
	return m.addCallback(time.Duration(0), n, callback, true)
}

// 加任务
func (m *minHeap) addCallback(expire time.Duration, n Next, callback func(), isSchedule bool) TimeNoder {
	select {
	case <-m.ctx.Done():
		panic("cannot add a task to a closed timer")
	default:
	}

	node := minHeapNode{
		callback:   callback,
		userExpire: expire,
		next:       n,
		absExpire:  time.Now().Add(expire),
		isSchedule: isSchedule,
		root:       m,
	}

	if n != nil {
		node.absExpire = n.Next(time.Now())
	}

	m.mu.Lock()
	heap.Push(&m.minHeaps, &node)
	m.wait.Add(1)
	m.mu.Unlock()

	select {
	case m.chAdd <- struct{}{}:
	default:
	}

	return &node
}

func (m *minHeap) removeTimeNode(node *minHeapNode) {
	m.mu.Lock()
	if node.index < 0 || node.index > int32(len(m.minHeaps)) || int32(len(m.minHeaps)) == 0 {
		m.mu.Unlock()
		return
	}

	heap.Remove(&m.minHeaps, int(node.index))
	m.wait.Done()
	m.mu.Unlock()
}

func (m *minHeap) resetTimeNode(node *minHeapNode, d time.Duration) {
	m.mu.Lock()
	node.userExpire = d
	node.absExpire = time.Now().Add(d)
	heap.Fix(&m.minHeaps, int(node.index))
	select {
	case m.chAdd <- struct{}{}:
	default:
	}
	m.mu.Unlock()
}

func (m *minHeap) getNewSleepTime() time.Duration {
	if m.minHeaps.Len() == 0 {
		return time.Hour
	}

	timeout := time.Until(m.minHeaps[0].absExpire)
	if timeout < 0 {
		timeout = 0
	}
	return timeout
}

func (m *minHeap) process() {
	for {
		m.mu.Lock()
		now := time.Now()
		// 如果堆中没有元素，就等待
		// 这时候设置一个相对长的时间，避免空转cpu
		if m.minHeaps.Len() == 0 {
			m.tm.Reset(time.Hour)
			m.mu.Unlock()
			return
		}

		for {
			// 取出最小堆的第一个元素
			first := m.minHeaps[0]

			// 时间未到直接过滤掉
			// 只是跳过最近的循环
			if !now.After(first.absExpire) {
				break
			}

			// 取出待执行的callback
			callback := first.callback
			// 如果是周期性任务
			if first.isSchedule {
				// 计算下次触发的绝对时间点
				first.absExpire = first.Next(now)
				// 修改下在堆中的位置
				heap.Fix(&m.minHeaps, int(first.index))
			} else {
				// 从堆中删除
				heap.Pop(&m.minHeaps)
				m.wait.Done()
			}

			// 正在运行的任务数加1
			atomic.AddInt32(&m.runCount, 1)
			go func() {
				callback()
				// 对正在运行的任务数减1
				atomic.AddInt32(&m.runCount, -1)
			}()

			// 如果堆中没有元素，就等待
			if m.minHeaps.Len() == 0 {
				m.tm.Reset(defaultTimeout)
				m.mu.Unlock()
				return
			}
		}

		// 取出第一个元素
		first := m.minHeaps[0]
		// 如果第一个元素的时间还没到，就计算下次触发的时间
		if time.Now().Before(first.absExpire) {
			to := m.getNewSleepTime()
			m.tm.Reset(to)
			// fmt.Printf("### now=%v, to = %v, m.minHeaps[0].absExpire = %v\n", time.Now(), to, m.minHeaps[0].absExpire)
			m.mu.Unlock()
			return
		}
		m.mu.Unlock()
	}
}

// 运行
// 为了避免空转cpu, 会等待一个chan, 只要AfterFunc或者ScheduleFunc被调用就会往这个chan里面写值
func (m *minHeap) Run() {
	m.tm = time.NewTimer(time.Hour)
	m.process()
	for {
		select {
		case <-m.tm.C:
			m.process()
		case <-m.chAdd:
			m.mu.Lock()
			// 极端情况，加完任务立即给删除了, 判断下当前堆中是否有元素
			if m.minHeaps.Len() > 0 {
				m.tm.Reset(m.getNewSleepTime())
			}
			m.mu.Unlock()
			// 进入事件循环，如果为空就会从事件循环里面退出
		case <-m.ctx.Done():
			// 等待所有任务结束
			m.wait.Wait()
			return
		}

	}
}

// 停止所有定时器
func (m *minHeap) Stop() {
	m.cancel()
}

func newMinHeap() (mh *minHeap) {
	mh = &minHeap{}
	heap.Init(&mh.minHeaps)
	mh.chAdd = make(chan struct{}, 1024)
	mh.ctx, mh.cancel = context.WithCancel(context.TODO())
	return
}
