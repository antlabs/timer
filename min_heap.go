package timer

import "time"

var _ Timer = (*minHeap)(nil)

type minHeap struct {
}

// 一次性定时器
func (m minHeap) AfterFunc(expire time.Duration, callback func()) TimeNoder {
	return nil
}

// 周期性定时器
func (m minHeap) ScheduleFunc(expire time.Duration, callback func()) TimeNoder {
	return nil
}

// 运行
func (m minHeap) Run() {

}

// 停止所有定时器
func (m minHeap) Stop() {

}

func newMinHeap() *minHeap {
	return &minHeap{}
}
