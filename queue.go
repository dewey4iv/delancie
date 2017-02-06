package delancie

// Queue defines a FIFO queue
type Queue interface {
	Push(b []byte) error
	Pop() ([]byte, error)
	Done()
	Stop()
}

// ChanQueue defines a FIFO queue that uses channels instead of functions
type ChanQueue interface {
	In() chan []byte
	Out() chan []byte
	Errors() chan error
	Done()
	Stop()
}
