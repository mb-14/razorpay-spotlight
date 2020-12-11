package main

import (
	"flag"
)

const (
	Realtime = "realtime"
	Backfill = "backfill"
)

const (
	Endpoint = "https://b9ef400928a48731d4888532d0493535.m.pipedream.net"
)

// Config
type Config struct {
	Event              string `yaml:"event"`
	MethodDistribution struct {
		Netbanking uint `yaml:"netbanking"`
		Card       uint `yaml:"card"`
		Upi        uint `yaml:"upi"`
		Wallet     uint `yaml:"wallet"`
	} `yaml:"method_distribution"`
	EnumFields []struct {
		Path   string   `yaml:"path"`
		Values []string `yaml:"values"`
	} `yaml:"enum_fields"`
	RangeFields []struct {
		Path string `yaml:"path"`
		Min  int    `yaml:"min"`
		Max  int    `yaml:"max"`
	} `yaml:"range_fields"`
}

var interval = flag.Int64("interval", 1000, "Number of milliseconds between events")
var event = flag.String("event", "payment_authorized", "Event type")
var mode = flag.String("mode", "realtime", "Mode: realtime or backfill")
var durationInDays = flag.Int("duration", 10, "Duration in days")
var config Config

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	flag.Parse()

	if *mode == Backfill {
		backfillEvents()
	}
}
