package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func backfillEvents() {
	timestamps := generateTimestamps()
	fmt.Printf("Number of events: %d\n", len(timestamps))
	numJobs := len(timestamps)
	jobs := make(chan time.Time, numJobs)
	results := make(chan error, numJobs)

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, t := range timestamps {
		jobs <- t
	}

	close(jobs)

	// Finally we collect all the results of the work.
	// This also ensures that the worker goroutines have
	// finished. An alternative way to wait for multiple
	// goroutines is to use a [WaitGroup](waitgroups).
	for a := 1; a <= numJobs; a++ {
		err := <-results
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func worker(id int, jobs <-chan time.Time, results chan<- error) {
	for j := range jobs {
		payload := bytes.NewBuffer(generatePayload(j))
		resp, err := http.Post(Endpoint, "application/json", payload)
		if resp.StatusCode != http.StatusOK {
			var message []byte
			resp.Body.Read(message)
			err = fmt.Errorf("%d : %s", resp.StatusCode, string(message))
		}
		results <- err
	}
}

func generateTimestamps() []time.Time {
	rand.Seed(time.Now().UTC().UnixNano())
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -*durationInDays)
	timestamps := []time.Time{startTime}
	for timestamps[len(timestamps)-1].Before(endTime) {
		randInterval := random(*interval, 0.2)
		newTimestamp := timestamps[len(timestamps)-1].Add(time.Millisecond * time.Duration(randInterval))
		timestamps = append(timestamps, newTimestamp)
	}
	return timestamps
}

func random(num int64, variance float64) int64 {
	deviation := int64(variance * float64(num))
	return rand.Int63n(2*deviation) + num - deviation
}
