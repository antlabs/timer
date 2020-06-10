package main

import (
	"log"
	"sync"
	"time"

	"github.com/antlabs/timer"
)

// 这是一个长时间测试代码

// 测试周期执行
func schedule(tm timer.Timer) {
	tm.ScheduleFunc(200*time.Millisecond, func() {
		log.Printf("schedule 200 milliseconds\n")
	})

	tm.ScheduleFunc(time.Second, func() {
		log.Printf("schedule second\n")
	})

	tm.ScheduleFunc(1*time.Minute, func() {
		log.Printf("schedule minute\n")
	})

	tm.ScheduleFunc(1*time.Hour, func() {
		log.Printf("schedule hour\n")
	})

	tm.ScheduleFunc(24*time.Hour, func() {
		log.Printf("schedule day\n")
	})
}

// 测试一次性定时器
func after(tm timer.Timer) {
	var wg sync.WaitGroup
	wg.Add(4)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for i := 0; i < 3600*24; i++ {
			i := i
			tm.AfterFunc(time.Second, func() {
				log.Printf("after second:%d\n", i)
			})
			time.Sleep(900 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 60*24; i++ {
			i := i
			tm.AfterFunc(time.Minute, func() {
				log.Printf("after minute:%d\n", i)
			})
			time.Sleep(50 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 24; i++ {
			i := i
			tm.AfterFunc(time.Hour, func() {
				log.Printf("after hour:%d\n", i)
			})
			time.Sleep(59 * time.Minute)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1; i++ {
			i := i
			tm.AfterFunc(24*time.Hour, func() {
				log.Printf("after day:%d\n", i)
			})
			time.Sleep(59 * time.Minute)
		}
	}()
}

// 检测 stop after 消息，没有打印是正确的行为
func stopNode(tm timer.Timer) {

	var wg sync.WaitGroup
	wg.Add(4)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for i := 0; i < 3600*24; i++ {
			i := i
			node := tm.AfterFunc(time.Second, func() {
				log.Printf("stop after second:%d\n", i)
			})
			time.Sleep(900 * time.Millisecond)
			node.Stop()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 60*24; i++ {
			i := i
			node := tm.AfterFunc(time.Minute, func() {
				log.Printf("stop after minute:%d\n", i)
			})
			time.Sleep(50 * time.Second)
			node.Stop()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 24; i++ {
			i := i
			node := tm.AfterFunc(time.Hour, func() {
				log.Printf("stop after hour:%d\n", i)
			})
			time.Sleep(59 * time.Minute)
			node.Stop()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1; i++ {
			i := i
			node := tm.AfterFunc(23*time.Hour, func() {
				log.Printf("stop after day:%d\n", i)
			})
			time.Sleep(22 * time.Hour)
			node.Stop()
		}
	}()
}

func main() {

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	tm := timer.NewTimer()

	go schedule(tm)
	go after(tm)
	go stopNode(tm)

	go func() {
		time.Sleep(time.Hour*24 + time.Hour)
		tm.Stop()
	}()
	tm.Run()
}
