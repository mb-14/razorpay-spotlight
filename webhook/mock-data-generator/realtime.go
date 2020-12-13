package main

import (
	"time"

	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mb-14/rzp-spotlight/webhook/rzp"
)

func generateRealtime() {
	client := influxdb2.NewClient(DBEndpoint, "")
	writeAPI := client.WriteAPIBlocking("", DBName)
	for {
		payload := generatePayloadJson(time.Now())
		point := rzp.ProcessPayloadJson(payload)
		writeAPI.WritePoint(context.Background(), point)
		randInterval := random(*interval, 0.35)
		time.Sleep(time.Millisecond * time.Duration(randInterval))
	}
}
