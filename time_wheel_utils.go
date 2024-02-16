// Copyright 2020-2024 guonaihong, antlabs. All rights reserved.
//
// mit license
package timer

import "time"

func get10Ms() time.Duration {
	return time.Duration(int64(time.Now().UnixNano() / int64(time.Millisecond) / 10))
}
