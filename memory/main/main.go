package main

import (
	"log"
	"time"
)

func main() {
	info("Starting...")

	q := New()
	ticker := NewTicker(time.Millisecond * 1)
	go func() {
		for now := range ticker.C {
			if now != nil {
				q.Push([]byte(now.Format("04:05.000")))
			} else {
				info("Done Creating.")
				q.Done()
				return
			}
		}
	}()

	go func() {
		time.Sleep(time.Millisecond * 1000 * 2)
		ticker.Stop()
	}()

	for {
		output := q.Pop()

		if output == "" {
			q.Stop()
			break
		}

		info(output)
		time.Sleep(time.Millisecond * 10)
	}

}

func info(args ...interface{}) {
	log.Println(args...)
}
