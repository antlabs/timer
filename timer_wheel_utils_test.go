package timer

import (
	"fmt"
	"testing"
)

func Test_Get10Ms(t *testing.T) {

	fmt.Printf("%v:%d", get10Ms(), get10Ms())
}
