package memory

import "errors"

// Option defines a thing that can be applied to a Queue
type Option interface {
	Apply(*Queue) error
}

// WithDefaults is automatically applied to a Queue and all others rewrite
// it's properties
func WithDefaults() Option {
	return &withDefaults{}
}

type withDefaults struct{}

func (opt *withDefaults) Apply(q *Queue) error {
	// default to a 10MB buffer - it's small, but safe
	if err := WithMaxBufferLen(1000).Apply(q); err != nil {
		return err
	}

	return nil
}

// WithMaxBufferLen sets the max number of bytes in a particular buffer
func WithMaxBufferLen(size int) Option {
	return &withMaxBufferLen{size}
}

type withMaxBufferLen struct {
	size int
}

func (opt *withMaxBufferLen) Apply(q *Queue) error {
	if opt.size < 1 {
		return errors.New("buffer size must be larger than 0")
	}

	q.maxLen = opt.size

	return nil
}
