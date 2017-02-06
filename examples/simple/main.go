package main

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"

	"github.com/dewey4iv/delancie"
	"github.com/dewey4iv/delancie/memory"
)

func main() {
	log.Println("Starting...")

	totalRecords := 1000001

	queue, err := memory.New()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(queue delancie.Queue) {
		// Generate a bunch of records
		var buf bytes.Buffer
		for record := range delancie.RecordChan(totalRecords) {
			if err := json.NewEncoder(&buf).Encode(record); err != nil {
				panic(err)
			}

			err := queue.Push(buf.Bytes())
			if err != nil {
				log.Printf("[ERROR]: %s", err.Error())
			}

			wg.Add(1)
			buf.Reset()
		}
	}(queue)

	var counter int
	for {
		b, err := queue.Pop()
		if err != nil {
			log.Println("[ERROR]:", err.Error())
			continue
		}

		wg.Done()
		counter++

		if true {
			log.Println("Read:", string(b))
		}

		if counter >= totalRecords {
			break
		}
	}

	wg.Done()
	wg.Wait()

	log.Println("End.")
}
