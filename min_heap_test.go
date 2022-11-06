package timer

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试AfterFunc有没有运行以及时间间隔可对
func Test_MinHeap_AfterFunc_Run(t *testing.T) {

	t.Run("1ms", func(t *testing.T) {

		tm := NewTimer(WithMinHeap())
		now := time.Now()
		go tm.Run()
		count := int32(0)

		tc := make(chan time.Duration, 2)

		tm.AfterFunc(time.Millisecond, func() {
			atomic.AddInt32(&count, 1)
			tc <- time.Since(now)
		})
		tm.AfterFunc(time.Millisecond, func() {
			atomic.AddInt32(&count, 1)
			tc <- time.Since(now)
		})

		time.Sleep(time.Millisecond * 2)
		close(tc)
		for tv := range tc {
			if tv < time.Millisecond || tv > 2*time.Millisecond {
				assert.Fail(t, "tc < time.Millisecond tc > 2*time.Millisecond")
			}
		}
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

	t.Run("10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		now := time.Now()
		go tm.Run()
		count := int32(0)
		tc := make(chan time.Duration, 2)
		tm.AfterFunc(time.Millisecond*10, func() {
			atomic.AddInt32(&count, 1)
			tc <- time.Since(now)
		})
		tm.AfterFunc(time.Millisecond*10, func() {
			atomic.AddInt32(&count, 1)
			tc <- time.Since(now)
		})

		time.Sleep(time.Millisecond * 20)
		close(tc)
		cnt := 1
		for tv := range tc {
			left := time.Millisecond * 10 * time.Duration(cnt)
			right := time.Duration(cnt) * 2 * 10 * time.Millisecond
			if tv < left || tv > right {
				t.Errorf("index(%d) (%v)tc < %v || tc > %v", cnt, tv, left, right)
			}
			//cnt++
		}
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

	t.Run("90ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		tm.AfterFunc(time.Millisecond*90, func() { atomic.AddInt32(&count, 1) })
		tm.AfterFunc(time.Millisecond*90, func() { atomic.AddInt32(&count, 1) })

		time.Sleep(time.Millisecond * 180)
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

}

// 测试Schedule 运行的周期可对
func Test_MinHeap_ScheduleFunc_Run(t *testing.T) {

	t.Run("1ms", func(t *testing.T) {

		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		c := make(chan bool, 1)
		node := tm.ScheduleFunc(time.Millisecond, func() {
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) == 2 {
				c <- true
			}
		})

		go func() {
			<-c
			node.Stop()
			node.Stop()
		}()

		time.Sleep(time.Millisecond * 5)
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

	t.Run("10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		c := make(chan bool, 1)
		tc := make(chan time.Duration, 2)
		now := time.Now()
		node := tm.ScheduleFunc(time.Millisecond*10, func() {
			atomic.AddInt32(&count, 1)
			tc <- time.Since(now)
			if atomic.LoadInt32(&count) == 2 {
				c <- true
			}
		})

		go func() {
			<-c
			node.Stop()
		}()

		time.Sleep(time.Millisecond * 30)
		close(tc)
		cnt := 1
		for tv := range tc {
			left := time.Millisecond * 10 * time.Duration(cnt)
			right := time.Duration(cnt) * 2 * 10 * time.Millisecond
			if tv < left || tv > right {
				t.Errorf("index(%d) (%v)tc < %v || tc > %v", cnt, tv, left, right)
			}
			cnt++
		}
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

	t.Run("30ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		c := make(chan bool, 1)

		node := tm.ScheduleFunc(time.Millisecond*30, func() {
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) == 2 {
				c <- true
			}
		})
		go func() {

			<-c
			node.Stop()
		}()

		time.Sleep(time.Millisecond * 70)
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})
}

// 测试Stop是否会等待正在运行的任务结束
func Test_Run_Stop(t *testing.T) {
	t.Run("1ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		count := uint32(0)
		tm.AfterFunc(time.Millisecond, func() { atomic.AddUint32(&count, 1) })
		tm.AfterFunc(time.Millisecond, func() { atomic.AddUint32(&count, 1) })
		go func() {
			time.Sleep(time.Millisecond * 4)
			tm.Stop()
		}()
		tm.Run()
		assert.Equal(t, atomic.LoadUint32(&count), uint32(2))
	})
}

type curstomTest struct {
	count int
}

func (c *curstomTest) Next(now time.Time) (rv time.Time) {
	rv = now.Add(time.Duration(c.count) * time.Millisecond * 10)
	c.count++
	return
}

func Test_CustomFunc(t *testing.T) {
	t.Run("custom", func(t *testing.T) {

		tm := NewTimer(WithMinHeap())
		mh := tm.(*minHeap)
		tc := make(chan time.Duration, 2)
		now := time.Now()
		count := uint32(1)
		node := tm.CustomFunc(&curstomTest{count: 1}, func() {
			if atomic.LoadUint32(&count) == 2 {
				return
			}
			atomic.AddUint32(&count, 1)
			tc <- time.Since(now)
		})

		go func() {
			time.Sleep(time.Millisecond * 30)
			node.Stop()
			tm.Stop()
		}()

		tm.Run()
		close(tc)
		cnt := 1
		for tv := range tc {
			left := time.Millisecond * 10 * time.Duration(cnt)
			right := time.Duration(cnt) * 2 * 10 * time.Millisecond
			if tv < left || tv > right {
				t.Errorf("index(%d) (%v)tc < %v || tc > %v", cnt, tv, left, right)
			}
			cnt++
		}
		assert.Equal(t, atomic.LoadUint32(&count), uint32(2))
		assert.Equal(t, mh.runCount, uint32(3))
	})
}
