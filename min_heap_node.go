// Copyright 2020-2024 guonaihong, antlabs. All rights reserved.
//
// mit license
package timer

import (
	"time"
)

type minHeapNode struct {
	callback   func()        // 用户的callback
	absExpire  time.Time     // 绝对时间
	userExpire time.Duration // 过期时间段
	root       *minHeap      // 指向最小堆
	next       Next          // 自定义下个触发的时间点, cronex项目用到了
	index      int32         // 在min heap中的索引，方便删除或者重新推入堆中
	isSchedule bool          // 是否是周期性任务
}

func (m *minHeapNode) Stop() bool {
	m.root.removeTimeNode(m)
	return true
}
func (m *minHeapNode) Reset(d time.Duration) bool {
	m.root.resetTimeNode(m, d)
	return true
}

func (m *minHeapNode) Next(now time.Time) time.Time {
	if m.next != nil {
		return (m.next).Next(now)
	}
	return now.Add(m.userExpire)
}

type minHeaps []*minHeapNode

func (m minHeaps) Len() int { return len(m) }

func (m minHeaps) Less(i, j int) bool { return m[i].absExpire.Before(m[j].absExpire) }

func (m minHeaps) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
	m[i].index = int32(i)
	m[j].index = int32(j)
}

func (m *minHeaps) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*m = append(*m, x.(*minHeapNode))
	lastIndex := int32(len(*m) - 1)
	(*m)[lastIndex].index = lastIndex
}

func (m *minHeaps) Pop() any {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}
