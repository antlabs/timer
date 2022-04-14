package timer

import "time"

// 定时器接口
type Timer interface {
	// 一次性定时器
	AfterFunc(expire time.Duration, callback func(args ...interface{}), args ...interface{}) TimeNoder

	// 周期性定时器
	ScheduleFunc(expire time.Duration, callback func(args ...interface{}), args ...interface{}) TimeNoder

	// 运行
	Run()

	// 停止所有定时器
	Stop()
}

// 停止单个定时器
type TimeNoder interface {
	Stop()
}

// 定时器构造函数
func NewTimer() Timer {
	return newTimeWheel()
}
