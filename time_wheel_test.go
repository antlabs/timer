package timer

import (
	"math"
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

func Test_hour(t *testing.T) {
	tw := newTimeWheel()

	testHour := new(bool)
	done := make(chan struct{}, 1)
	tw.AfterFunc(time.Hour, func() {
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
