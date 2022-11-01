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
	nodeHeaps
	chAdd  chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

// 一次性定时器
func (m *minHeap) AfterFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, callback, false)
}

func (m *minHeap) addCallback(expire time.Duration, callback func(), isSchedule bool) TimeNoder {
	m.mu.Lock()
	defer m.mu.Unlock()

	node := minHeapNode{
		callback:   callback,
		userExpire: expire,
	}
	heap.Push(&m.nodeHeaps, node)
	return &node

}

// 周期性定时器
func (m *minHeap) ScheduleFunc(expire time.Duration, callback func()) TimeNoder {
	return m.addCallback(expire, callback, true)
}

// 运行
// 为了避免空转cpu, 会等待一个chan, 只要AfterFunc或者SchedulerFunc被调用就会往这个chan里面写值
func (m *minHeap) Run() {
	for {
		select {
		case <-m.chAdd:
			// 进入事件循环，如果为空就会从事件循环里面退出
			for {

			}
		case <-m.ctx.Done():
			// 等待所有任务结束
			//return
		}
	}
}

// 停止所有定时器
func (m *minHeap) Stop() {

}

func newMinHeap() (mh *minHeap) {
	mh = &minHeap{}
	heap.Init(&mh.nodeHeaps)
	mh.chAdd = make(chan struct{})
	mh.ctx, mh.cancel = context.WithCancel(context.TODO())
	return
}
