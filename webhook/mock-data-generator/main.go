package main

import (
	"flag"
)

const (
	Realtime = "realtime"
	Backfill = "backfill"
)

const (
	Endpoint = "http://159.65.149.215:8080/webhook"
)

var interval = flag.Int64("interval", 1000, "Number of milliseconds between events")
var event = flag.String("event", "payment_authorized", "Event type")
var mode = flag.String("mode", "realtime", "Mode: realtime or backfill")
var durationInDays = flag.Int("duration", 10, "Duration in days")

func main() {
	flag.Parse()
	if *mode == Backfill {
		backfillEvents()
	}
}
