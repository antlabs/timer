package timer

import (
	"testing"
	"time"
)

func TestHelloWorld(t *testing.T) {
	tm := &timer{}
	tm.add(time.Millisecond*20, nil)
	tm.add(time.Second, nil)
	tm.add(time.Minute, nil)
	tm.add(time.Hour, nil)
	tm.add(time.Hour*24, nil)
	tm.add(time.Hour*24*365, nil)
	tm.add(time.Hour*24*365*12, nil)
	//tm.debug()
}
