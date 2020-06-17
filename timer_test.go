package timer

import (
	"log"
	"sync/atomic"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
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
	assert.Equal(t, atomic.LoadUint32(&count), uint32(5))
}

func Test_AfterFunc(t *testing.T) {
	tm := NewTimer()

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

	go func() {
		time.Sleep(time.Second + time.Millisecond*100)
		tm.Stop()
	}()
	tm.Run()
	assert.Equal(t, atomic.LoadUint32(&count), uint32(2))
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
	assert.NotEqual(t, count, 1)
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

	assert.NotEqual(t, count, 1)
}
