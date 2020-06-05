package timer

import (
	"log"
	"testing"
	"time"
)

func TestHelloWorld(t *testing.T) {
	tm := NewTimer()

	log.Printf("start\n")

	tm.AfterFunc(time.Millisecond*20, func() {
		log.Printf("after Millisecond * 20")
	})

	tm.AfterFunc(time.Second, func() {
		log.Printf("after second")
	})

	tm.AfterFunc(time.Minute, func() {
		log.Printf("after Minute")
	})

	/*
		tm.AfterFunc(time.Hour, nil)
		tm.AfterFunc(time.Hour*24, nil)
		tm.AfterFunc(time.Hour*24*365, nil)
		tm.AfterFunc(time.Hour*24*365*12, nil)
	*/
	go func() {
		time.Sleep(time.Minute + time.Second*4)
		tm.Stop()
	}()
	tm.Run()
	//tm.debug()
}
