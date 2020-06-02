package timer

import "time"

func getMs() time.Duration {
	return int64(time.Now().UnixNano() / 1000000)
}
