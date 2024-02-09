package timer

import (
	"container/heap"
	"testing"
	"time"
)

func Test_NodeSizeof(t *testing.T) {
	t.Run("输出最小堆node的sizeof", func(t *testing.T) {
		// t.Logf("minHeapNode size: %d， %d\n", unsafe.Sizeof(minHeapNode{}), unsafe.Sizeof(time.Timer{}))
	})
}
func Test_MinHeap(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var mh minHeaps
		now := time.Now()
		n1 := minHeapNode{
			absExpire:  now.Add(time.Second),
			userExpire: 1 * time.Second,
		}

		n2 := minHeapNode{
			absExpire:  now.Add(2 * time.Second),
			userExpire: 2 * time.Second,
		}

		n3 := minHeapNode{
			absExpire:  now.Add(3 * time.Second),
			userExpire: 3 * time.Second,
		}

		n6 := minHeapNode{
			absExpire:  now.Add(6 * time.Second),
			userExpire: 6 * time.Second,
		}
		n5 := minHeapNode{
			absExpire:  now.Add(5 * time.Second),
			userExpire: 5 * time.Second,
		}
		n4 := minHeapNode{
			absExpire:  now.Add(4 * time.Second),
			userExpire: 4 * time.Second,
		}
		mh.Push(&n1)
		mh.Push(&n2)
		mh.Push(&n3)
		mh.Push(&n6)
		mh.Push(&n5)
		mh.Push(&n4)

		for i := 1; len(mh) > 0; i++ {
			v := heap.Pop(&mh).(*minHeapNode)

			if v.userExpire != time.Duration(i)*time.Second {
				t.Errorf("index(%d) v.userExpire(%v) != %v", i, v.userExpire, time.Duration(i)*time.Second)
			}
		}
	})
}
