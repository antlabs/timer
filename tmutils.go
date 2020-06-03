package timer

import "time"

func get10Ms() time.Duration {
	return int64(time.Now().UnixNano() / 1000 / 1000 / 10)
}
