package delancie

import (
	"bytes"
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TestRecord is an exported struct for testing
type TestRecord struct {
	ID      string    `json:"id"`
	Value   int       `json:"value"`
	Created time.Time `json:"created"`
}

// RecordChan returns a channel that spits out a number of records equal to the provided int
func RecordChan(recordNum int) chan TestRecord {
	ch := make(chan TestRecord)

	go func() {
		for i := 0; i < recordNum; i++ {
			ch <- TestRecord{uuid.NewV4().String(), i, time.Now()}
			time.Sleep(time.Millisecond)
		}

		close(ch)
	}()

	return ch
}

// GenerateJSON creates a json stream with the length of recordNum. Primarily used for testing
func GenerateJSON(recordNum int) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	for i := 0; i < recordNum; i++ {
		if err := encoder.Encode(TestRecord{uuid.NewV4().String(), i, time.Now()}); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
