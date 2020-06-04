package timer

import (
	"container/list"
	"sync"
)

// 先使用sync.Mutex实现功能
// 后面使用cas优化
type Time struct {
	*list.List

	sync.Mutex
}

type timeNode struct {
	expire   uint64
	callback func()

	list *Time
}

func (t *timeNode) Remove() {
	t.list.Lock()
	defer t.list.Unlock()
	//t.list.Remove(t)
}

/*
func (t *timeNode) Stop() {
}
*/
