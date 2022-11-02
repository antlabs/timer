package timer

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MinHeap_AfterFunc(t *testing.T) {

	t.Run("1ms", func(t *testing.T) {

		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		tm.AfterFunc(time.Millisecond, func() { atomic.AddInt32(&count, 1) })
		tm.AfterFunc(time.Millisecond, func() { atomic.AddInt32(&count, 1) })

		time.Sleep(time.Millisecond * 2)
		assert.Equal(t, atomic.LoadInt32(&count), int32(2))

	})

	t.Run("10ms", func(t *testing.T) {
		tm := NewTimer(WithMinHeap())
		go tm.Run()
		count := int32(0)
		tm.AfterFunc(time.Millisecond*10, func() { atomic.AddInt32(&count, 1) })
		tm.AfterFunc(time.Millisecond*10, func() { atomic.AddInt32(&count, 1) })

		time.Sleep(time.Millisecond * 20)
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
