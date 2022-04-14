package timer

import (
	"context"
	"math"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_maxVal(t *testing.T) {
	assert.Equal(t, maxVal(), uint64(math.MaxUint32))
}

func Test_LevelMax(t *testing.T) {
	assert.Equal(t, levelMax(1), uint64(1<<(nearShift+levelShift)))
	assert.Equal(t, levelMax(2), uint64(1<<(nearShift+2*levelShift)))
	assert.Equal(t, levelMax(3), uint64(1<<(nearShift+3*levelShift)))
	assert.Equal(t, levelMax(4), uint64(1<<(nearShift+4*levelShift)))
}

func Test_GenVersion(t *testing.T) {
	assert.Equal(t, genVersionHeight(1, 0xf), uint64(0x0001000f00000000))
	assert.Equal(t, genVersionHeight(1, 64), uint64(0x0001004000000000))
}

// 测试1小时
func Test_hour(t *testing.T) {
	tw := newTimeWheel()

	testHour := new(bool)
	done := make(chan struct{}, 1)
	tw.AfterFunc(time.Hour, func(args ...interface{}) {
		*testHour = true
		done <- struct{}{}
	})

	expire := getExpire(time.Hour, 0)
	for i := 0; i < int(expire)+10; i++ {
		get10Ms := func() time.Duration {
			return tw.curTimePoint + 1
		}
		tw.run(get10Ms)
	}

	select {
	case <-done:
	case <-time.After(time.Second / 100):
	}
	assert.True(t, *testHour)
}

// 测试周期性定时器, 5s
func Test_ScheduleFunc_5s(t *testing.T) {
	tw := newTimeWheel()

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(1000)

	testTime := time.Second * 5

	tw.ScheduleFunc(testTime, func(args ...interface{}) {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		get10Ms := func() time.Duration {
			return tw.curTimePoint + 1
		}
		tw.run(get10Ms)
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)

}

// 测试周期性定时器, 1hour
func Test_ScheduleFunc_hour(t *testing.T) {
	tw := newTimeWheel()

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(100)
	testTime := time.Hour

	tw.ScheduleFunc(testTime, func(args ...interface{}) {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		get10Ms := func() time.Duration {
			return tw.curTimePoint + 1
		}
		tw.run(get10Ms)
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)

}

// 测试周期性定时器, 1day
func Test_ScheduleFunc_day(t *testing.T) {
	tw := newTimeWheel()

	var first5 int32
	ctx, cancel := context.WithCancel(context.Background())

	const total = int32(10)
	testTime := time.Hour * 24

	tw.ScheduleFunc(testTime, func(args ...interface{}) {
		atomic.AddInt32(&first5, 1)
		if atomic.LoadInt32(&first5) == total {
			cancel()
		}

	})

	expire := getExpire(testTime*time.Duration(total), 0)
	for i := 0; i <= int(expire)+10; i++ {
		get10Ms := func() time.Duration {
			return tw.curTimePoint + 1
		}
		tw.run(get10Ms)
	}

	select {
	case <-ctx.Done():
	case <-time.After(time.Second / 100):
	}

	assert.Equal(t, total, first5)

}
