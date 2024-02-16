// Copyright 2020-2024 guonaihong, antlabs. All rights reserved.
//
// mit license
package timer

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_ScheduleFunc(t *testing.T) {
	tm := NewTimer()

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	count := uint32(0)
	log.Printf("start\n")

	tm.ScheduleFunc(time.Millisecond*100, func() {
		log.Printf("schedule\n")
		atomic.AddUint32(&count, 1)
	})

	go func() {
		time.Sleep(570 * time.Millisecond)
		log.Printf("stop\n")
		tm.Stop()
	}()

	tm.Run()
	if count != 5 {
		t.Errorf("count:%d != 5\n", count)
	}

}

func Test_AfterFunc(t *testing.T) {
	tm := NewTimer()
	go tm.Run()
	log.Printf("start\n")

	count := uint32(0)
	tm.AfterFunc(time.Millisecond*20, func() {
		log.Printf("after Millisecond * 20")
		atomic.AddUint32(&count, 1)
	})

	tm.AfterFunc(time.Second, func() {
		log.Printf("after second")
		atomic.AddUint32(&count, 1)
	})

	/*
		tm.AfterFunc(time.Minute, func() {
			log.Printf("after Minute")
		})
	*/
	/*
		tm.AfterFunc(time.Hour, nil)
		tm.AfterFunc(time.Hour*24, nil)
		tm.AfterFunc(time.Hour*24*365, nil)
		tm.AfterFunc(time.Hour*24*365*12, nil)
	*/

	time.Sleep(time.Second + time.Millisecond*100)
	tm.Stop()

	if count != 2 {
		t.Errorf("count:%d != 2\n", count)
	}

}

func Test_Node_Stop_1(t *testing.T) {
	tm := NewTimer()
	count := uint32(0)
	node := tm.AfterFunc(time.Millisecond*10, func() {
		atomic.AddUint32(&count, 1)
	})
	go func() {
		time.Sleep(time.Millisecond * 30)
		node.Stop()
		tm.Stop()
	}()

	tm.Run()
	if count != 1 {
		t.Errorf("count:%d == 1\n", count)
	}
}

func Test_Node_Stop(t *testing.T) {
	tm := NewTimer()
	count := uint32(0)
	node := tm.AfterFunc(time.Millisecond*100, func() {
		atomic.AddUint32(&count, 1)
	})
	node.Stop()
	go func() {
		time.Sleep(time.Millisecond * 200)
		tm.Stop()
	}()
	tm.Run()

	if count == 1 {
		t.Errorf("count:%d == 1\n", count)
	}

}

// 测试重置定时器
func Test_Reset(t *testing.T) {
	t.Run("min heap reset", func(t *testing.T) {

		tm := NewTimer(WithMinHeap())

		go tm.Run()
		count := int32(0)

		tc := make(chan time.Duration, 2)

		var mu sync.Mutex
		isClose := false
		now := time.Now()
		node1 := tm.AfterFunc(time.Millisecond*100, func() {

			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		node2 := tm.AfterFunc(time.Millisecond*100, func() {
			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})
		node1.Reset(time.Millisecond)
		node2.Reset(time.Millisecond)

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

	t.Run("time wheel reset", func(t *testing.T) {
		tm := NewTimer()

		go func() {
			tm.Run()
		}()

		count := int32(0)

		tc := make(chan time.Duration, 2)

		var mu sync.Mutex
		isClose := false
		now := time.Now()
		node1 := tm.AfterFunc(time.Millisecond*10, func() {

			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		node2 := tm.AfterFunc(time.Millisecond*10, func() {
			mu.Lock()
			atomic.AddInt32(&count, 1)
			if atomic.LoadInt32(&count) <= 2 && !isClose {
				tc <- time.Since(now)
			}
			mu.Unlock()
		})

		node1.Reset(time.Millisecond * 20)
		node2.Reset(time.Millisecond * 20)

		time.Sleep(time.Millisecond * 40)
		mu.Lock()
		isClose = true
		close(tc)
		node1.Stop()
		node2.Stop()
		mu.Unlock()
		for tv := range tc {
			if tv < time.Millisecond*20 || tv > 2*time.Millisecond*20 {
				t.Errorf("tc < time.Millisecond tc > 2*time.Millisecond")
			}
		}
		if atomic.LoadInt32(&count) != 2 {
			t.Errorf("count:%d != 2", atomic.LoadInt32(&count))
		}
	})
}
