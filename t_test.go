package timer

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_Look(t *testing.T) {

	tmp := newTimeHead(0, 0)
	offset := unsafe.Offsetof(tmp.Head)
	fmt.Printf("%d\n", offset)
}
