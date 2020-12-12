package main

import (
	"flag"
)

const (
	Realtime          = "realtime"
	Backfill          = "backfill"
	PaymentAuthorized = "payment_authorized"
	PaymentFailed     = "payment_failed"
)

const (
	Endpoint = "http://159.65.149.215:8080/webhook"
	//Endpoint = "https://b9ef400928a48731d4888532d0493535.m.pipedream.net"
)

type MethodConfig struct {
	Weight uint `yaml:"weight"`
	Fields []struct {
		Path   string   `yaml:"path"`
		Values []string `yaml:"values"`
	} `yaml:"fields"`
}

type Config struct {
	Netbanking  MethodConfig `yaml:"netbanking"`
	Card        MethodConfig `yaml:"card"`
	Wallet      MethodConfig `yaml:"wallet"`
	Upi         MethodConfig `yaml:"upi"`
	RangeFields []struct {
		Path string `yaml:"path"`
		Min  int    `yaml:"min"`
		Max  int    `yaml:"max"`
	} `yaml:"range_fields"`
	ErrorFields []struct {
		Path   string   `yaml:"path"`
		Values []string `yaml:"values"`
	} `yaml:"error_fields"`
}

var interval = flag.Int64("interval", 5000, "Number of milliseconds between events")
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
