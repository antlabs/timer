## timer
[![Go](https://github.com/antlabs/timer/workflows/Go/badge.svg)](https://github.com/antlabs/timer/actions)
[![codecov](https://codecov.io/gh/antlabs/timer/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/timer)

timer是高性能定时器库
## feature
* 支持一次性定时器
* 支持周期性定时器

## 一次性定时器
```go
import (
    "github.com/antlabs/timer"
    "log"
)

func main() {
        tm := timer.NewTimer()

        tm.AfterFunc(1*time.Second, func() {
                log.Printf("after\n")
        })

        tm.AfterFunc(10*time.Second, func() {
                log.Printf("after\n")
        })
        tm.Run()
}
```
## 周期性定时器
```go
import (
    "github.com/antlabs/timer"
    "log"
)

func main() {
        tm := timer.NewTimer()

        tm.ScheduleFunc(1*time.Second, func() {
                log.Printf("schedule\n")
        })

        tm.Run()
}
```
