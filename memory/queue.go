package memory

import (
	"sync"
	"time"
)

// New takes a set of options and creates a new Queue
func New(opts ...Option) (*Queue, error) {
	q := Queue{
		100,
		sync.Mutex{},
		0,
		0,
		0,
		make(chan string),
		make(chan struct{}),
		make(chan string),
		make(chan struct{}),
		make([][]string, 1),
	}

	if opts == nil {
		if err := WithDefaults().Apply(&q); err != nil {
			return nil, err
		}
	} else {
		for _, opt := range opts {
			if err := opt.Apply(&q); err != nil {
				return nil, err
			}
		}
	}

	go q.loop()

	return &q, nil
}

type Queue struct {
	maxLen    int
	mux       sync.Mutex
	reader    int
	readerPos int
	writer    int
	pushCh    chan string
	popReqCh  chan struct{}
	popRespCh chan string
	closeCh   chan struct{}
	queue     [][]string
}

func (q *Queue) Push(b []byte) error {
	q.pushCh <- string(b)

	return nil
}

func (q *Queue) Pop() ([]byte, error) {
	q.popReqCh <- struct{}{}

	select {
	case data := <-q.popRespCh:
		return []byte(data), nil
	case <-q.closeCh:
		return nil, nil
	}
}

func (q *Queue) Done() {
	q.pushCh = nil
}

func (q *Queue) Stop() {
	close(q.closeCh)
}

func (q *Queue) loop() {
	for {
		select {
		case s := <-q.pushCh:
			q.writeOne(s)
		case <-q.popReqCh:
			if q.pushCh == nil {
				if q.reader == q.writer && q.readerPos == len(q.queue[q.reader]) {
					q.popRespCh <- ""
					return
				}
			}

			go func() {
				ticker := NewTicker(time.Millisecond)

				for {
					select {
					case <-ticker.C:
						if data := q.readOne(); data != "" {
							q.popRespCh <- data
							ticker.Stop()
							return
						}
					case <-q.closeCh:
						ticker.Stop()
						return
					}
				}
			}()
		case <-q.closeCh:
			return
		}
	}
}

func (q *Queue) writeOne(s string) {
	// check that we won't overflow
	if len(q.queue[q.writer]) >= q.maxLen {
		q.advanceWriter()
	}

	// add to queue
	q.queue[q.writer] = append(q.queue[q.writer], s)
}

func (q *Queue) readOne() string {
	q.mux.Lock()
	defer q.mux.Unlock()

	// if there are still spots in current reader
	if q.readerPos < len(q.queue[q.reader]) {
		q.readerPos++
		return q.queue[q.reader][q.readerPos-1]
	}

	// if reader has reached end
	if q.readerPos >= len(q.queue[q.reader]) {
		if q.reader < q.writer {
			q.advanceReader()
			q.readerPos++
			return q.queue[q.reader][q.readerPos-1]
		}

		if q.reader == q.writer {
			return ""
		}
	}

	panic("::: we shouldn't get here :::")
}

func (q *Queue) advanceReader() {
	q.queue[q.reader] = nil
	q.reader++
	q.readerPos = 0
}

func (q *Queue) advanceWriter() {
	q.queue = append(q.queue, []string{})
	q.writer++
}
