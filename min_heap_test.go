package timer

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 测试AfterFunc有没有运行以及时间间隔可对
func Test_MinHeap_AfterFunc_Run(t *testing.T) {
	t.Run("1ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())

		go tm.Run()
		count := int32(0)

		tc := make(chan time.Duration, 2)

		var mu sync.Mutex
		isClose := false
		now := time.Now()
		node1 := tm.AfterFunc(time.Millisecond, func() {

			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		node2 := tm.AfterFunc(time.Millisecond, func() {
			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		time.Sleep(time.Millisecond * 3)
		mu.Lock()
		isClose = true
		close(tc)
		node1.Stop()
		node2.Stop()
		mu.Unlock()
		for tv := range tc {
			if tv < time.Millisecond || tv > 2*time.Millisecond {
				t.Errorf("tc < time.Millisecond tc > 2*time.Millisecond")

			}
		}
		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}

	})

	t.Run("10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())

		go tm.Run() // 运行事件循环
		count := int32(0)
		tc := make(chan time.Duration, 2)

		var mu sync.Mutex
		isClosed := false
		now := time.Now()
		node1 := tm.AfterFunc(time.Millisecond*10, func() {
			now2 := time.Now()
			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClosed {
				tc <- time.Since(now)
			}
			mu.Unlock()
			log.Printf("node1.Lock:%v\n", time.Since(now2))
		})
		node2 := tm.AfterFunc(time.Millisecond*10, func() {
			now2 := time.Now()
			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClosed {
				tc <- time.Since(now)
			}
			mu.Unlock()
			log.Printf("node2.Lock:%v\n", time.Since(now2))
		})

		time.Sleep(time.Millisecond * 24)
		now3 := time.Now()
		mu.Lock()
		node1.Stop()
		node2.Stop()
		isClosed = true
		close(tc)
		mu.Unlock()

		log.Printf("node1.Stop:%v\n", time.Since(now3))
		cnt := 1
		for tv := range tc {
			left := time.Millisecond * 10 * time.Duration(cnt)
			right := time.Duration(cnt) * 2 * 10 * time.Millisecond
			if tv < left || tv > right {
				t.Errorf("index(%d) (%v)tc < %v || tc > %v", cnt, tv, left, right)
			}
			// cnt++
		}
		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}

	})

	t.Run("90ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		tm.AfterFunc(time.Millisecond*90, func() { atomic.AddInt32(&count, 1) })
		tm.AfterFunc(time.Millisecond*90, func() { atomic.AddInt32(&count, 2) })

		time.Sleep(time.Millisecond * 180)
		if atomic.LoadInt32(&count) != 3 {
			t.Errorf("count != 3")
		}

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
		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}

	})

	t.Run("10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		tc := make(chan time.Duration, 2)
		var mu sync.Mutex
		isClosed := false
		now := time.Now()

		node := tm.ScheduleFunc(time.Millisecond*10, func() {
			mu.Lock()
			atomic.AddInt32(&count, 1)

			if atomic.LoadInt32(&count) <= 2 && !isClosed {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		time.Sleep(time.Millisecond * 25)

		mu.Lock()
		close(tc)
		isClosed = true
		node.Stop()
		mu.Unlock()

		cnt := 1
		for tv := range tc {
			left := time.Millisecond * 10 * time.Duration(cnt)
			right := time.Duration(cnt) * 2 * 10 * time.Millisecond
			if tv < left || tv > right {
				t.Errorf("index(%d) (%v)tc < %v || tc > %v", cnt, tv, left, right)
			}
			cnt++
		}

		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}

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
		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}

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
		if atomic.LoadUint32(&count) != 2 {
			t.Errorf("count != 2")
		}
	})
}

type curstomTest struct {
	count int32
}

func (c *curstomTest) Next(now time.Time) (rv time.Time) {
	rv = now.Add(time.Duration(c.count) * time.Millisecond * 10)
	atomic.AddInt32(&c.count, 1)
	return
}

// 验证自定义函数的运行间隔时间
func Test_CustomFunc(t *testing.T) {
	t.Run("custom", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		// mh := tm.(*minHeap) // 最小堆
		tc := make(chan time.Duration, 2)
		now := time.Now()
		count := uint32(1)
		stop := make(chan bool, 1)
		// 自定义函数
		node := tm.CustomFunc(&curstomTest{count: 1}, func() {

			if atomic.LoadUint32(&count) == 2 {
				return
			}
			// 计算运行次数
			atomic.AddUint32(&count, 1)
			tc <- time.Since(now)
			// 关闭这个任务
			close(stop)
		})

		go func() {
			<-stop
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
		if atomic.LoadUint32(&count) != 2 {
			t.Errorf("count != 2")
		}

		// 正在运行的任务是比较短暂的，所以外部很难
		// if mh.runCount != int32(1) {
		// 	t.Errorf("mh.runCount:%d != 1", mh.runCount)
		// }

	})
}

// 验证运行次数是符合预期的
func Test_RunCount(t *testing.T) {
	t.Run("runcount-10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		max := 10
		go func() {
			tm.Run()
		}()

		count := uint32(0)
		for i := 0; i < max; i++ {
			tm.ScheduleFunc(time.Millisecond*10, func() {
				atomic.AddUint32(&count, 1)
			})
		}

		time.Sleep(time.Millisecond * 15)
		tm.Stop()
		if count != uint32(max) {
			t.Errorf("count != %d", max)
		}

	})
}
