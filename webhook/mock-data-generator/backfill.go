package main

import (
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mb-14/rzp-spotlight/webhook/rzp"
)

func backfillEvents() {
	client := influxdb2.NewClient(DBEndpoint, "")
	writeAPI := client.WriteAPI("", DBName)
	timestamps := generateTimestamps()
	for _, t := range timestamps {
		payload := generatePayloadJson(t)
		point := rzp.ProcessPayloadJson(payload)
		writeAPI.WritePoint(point)
	}
	client.Close()
}

func generateTimestamps() []time.Time {
	rand.Seed(time.Now().UTC().UnixNano())
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -*durationInDays)
	timestamps := []time.Time{startTime}
	for timestamps[len(timestamps)-1].Before(endTime) {
		randInterval := random(*interval, 0.35)
		newTimestamp := timestamps[len(timestamps)-1].Add(time.Millisecond * time.Duration(randInterval))
		timestamps = append(timestamps, newTimestamp)
	}
	return timestamps
}

func random(num int64, variance float64) int64 {
	deviation := int64(variance * float64(num))
	return rand.Int63n(2*deviation) + num - deviation
}
