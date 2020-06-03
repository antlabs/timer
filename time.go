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
}

func (t *timeNode) Remove() {
}

func (t *timeNode) Stop() {
}
