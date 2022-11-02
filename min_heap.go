package timer

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

var _ Timer = (*minHeap)(nil)

type minHeap struct {
	mu sync.Mutex
	minHeaps
	chAdd  chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wait   sync.WaitGroup
}

// 一次性定时器
func (m *minHeap) AfterFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, callback, false)
}

// 加任务
func (m *minHeap) addCallback(expire time.Duration, callback func(), isSchedule bool) TimeNoder {
	m.mu.Lock()
	defer m.mu.Unlock()

	node := minHeapNode{
		callback:   callback,
		userExpire: expire,
		absExpire:  time.Now().Add(expire),
		isSchedule: isSchedule,
	}

	heap.Push(&m.minHeaps, node)
	select {
	case m.chAdd <- struct{}{}:
	default:
	}

	return &node

}

// 周期性定时器
func (m *minHeap) ScheduleFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, callback, true)
}

// 运行
// 为了避免空转cpu, 会等待一个chan, 只要AfterFunc或者SchedulerFunc被调用就会往这个chan里面写值
func (m *minHeap) Run() {

	tm := time.NewTimer(time.Hour)
	for {
		select {
		case <-tm.C:
			for {
				m.mu.Lock()
				now := time.Now()
				if m.minHeaps.Len() == 0 {
					m.mu.Unlock()
					goto next
				}

				first := &m.minHeaps[0]
				var callback func()

				if now.After(first.absExpire) {
					if first.isSchedule {
						first.absExpire = now.Add(first.userExpire)
						heap.Fix(&m.minHeaps, first.index)
					} else {
						m.minHeaps.Pop()
					}
				}

				first = &m.minHeaps[0]
				if now.Before(first.absExpire) {
					tm.Reset(time.Since(m.minHeaps[0].absExpire))
				}
				m.mu.Unlock()
				if callback != nil {
					go callback()
				}
			}
		case <-m.chAdd:
			m.mu.Lock()
			// 极端情况，加完任务立即给删除了, 判断下当前堆中是否有元素
			if m.minHeaps.Len() > 0 {
				tm.Reset(time.Since(m.minHeaps[0].absExpire))
			}
			m.mu.Unlock()
			// 进入事件循环，如果为空就会从事件循环里面退出
		case <-m.ctx.Done():
			// 等待所有任务结束
			m.wait.Wait()
			return
		}
	next:
	}
}

// 停止所有定时器
func (m *minHeap) Stop() {
	m.cancel()
}

func newMinHeap() (mh *minHeap) {
	mh = &minHeap{}
	heap.Init(&mh.minHeaps)
	mh.chAdd = make(chan struct{}, 1)
	mh.ctx, mh.cancel = context.WithCancel(context.TODO())
	return
}
