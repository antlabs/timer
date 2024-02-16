// Copyright 2020-2024 guonaihong, antlabs. All rights reserved.
//
// mit license
package timer

import (
	"fmt"
	"testing"
)

func Test_Get10Ms(t *testing.T) {

	fmt.Printf("%v:%d", get10Ms(), get10Ms())
}
