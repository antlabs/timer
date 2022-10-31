package timer

type option struct {
	timeWheel bool
	minHeap   bool
}

type Option func(c *option)

func WithTimeWheel() Option {
	return func(o *option) {
		o.timeWheel = true
	}
}

func WithMinHeap() Option {
	return func(o *option) {
		o.minHeap = true
	}
}
