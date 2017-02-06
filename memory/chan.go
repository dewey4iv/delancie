package memory

func NewChan(opts ...Option) (*ChanQueue, error) {

	queue, err := New(opts...)
	if err != nil {
		return nil, err
	}

	chQ := ChanQueue{
		queue,
		make(chan []byte),
		make(chan []byte),
		make(chan error),
		make(chan struct{}),
	}

	go chQ.inLoop()
	go chQ.outLoop()

	return &chQ, nil
}

type ChanQueue struct {
	*Queue
	in      chan []byte
	out     chan []byte
	errors  chan error
	closeCh chan struct{}
}

func (q *ChanQueue) inLoop() {
	for {
		select {
		case input := <-q.in:
			q.Queue.Push(input)
		}
	}
}

func (q *ChanQueue) outLoop() {
	for {
		val, err := q.Pop()
		if err == nil {
			q.out <- val
		}
	}
}

func (q *ChanQueue) In() chan []byte {
	return q.in
}

func (q *ChanQueue) Out() chan []byte {
	return q.out
}

func (q *ChanQueue) Errors() chan error {
	return nil
}

func (q *ChanQueue) Stop() {
	q.Queue.Stop()
	close(q.closeCh)
}

func (q *ChanQueue) Done() {
	q.Queue.Done()
	close(q.closeCh)
}
