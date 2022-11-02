package timer

import (
	"time"
)

type minHeapNode struct {
	stop       uint32        //stop标记
	callback   func()        //用户的callback
	absExpire  time.Time     //绝对时间
	userExpire time.Duration //过期时间
	isSchedule bool          //是否是周期性任务
	index      int           //在min heap中的索引，方便删除用的
}

// TODO 实现
func (m *minHeapNode) Stop() {

}

type minHeaps []minHeapNode

func (m minHeaps) Len() int           { return len(m) }
func (m minHeaps) Less(i, j int) bool { return m[i].absExpire.Before(m[j].absExpire) }
func (m minHeaps) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m *minHeaps) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*m = append(*m, x.(minHeapNode))
	lastIndex := len(*m) - 1
	(*m)[lastIndex].index = lastIndex
}

func (m *minHeaps) Pop() any {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}
