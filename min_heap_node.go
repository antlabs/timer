package timer

import (
	"time"
)

type minHeapNode struct {
	stop       uint32        //stop标记
	callback   func()        //用户的callback
	userExpire time.Duration //过期时间
	isSchedule bool          //是否是周期性任务
	index      int           //在min heap中的索引，方便删除用的
}

// TODO 实现
func (m *minHeapNode) Stop() {

}

type nodeHeaps []minHeapNode

func (n nodeHeaps) Len() int           { return len(n) }
func (n nodeHeaps) Less(i, j int) bool { return n[i].userExpire < n[j].userExpire }
func (n nodeHeaps) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

func (n *nodeHeaps) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*n = append(*n, x.(minHeapNode))
	lastIndex := len(*n) - 1
	(*n)[lastIndex].index = lastIndex
}

func (h *nodeHeaps) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
