package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/mb-14/rzp-spotlight/webhook/json"
	"gopkg.in/yaml.v2"

	wr "github.com/mroth/weightedrand"
)

var methodChooser *wr.Chooser
var templates map[string]json.Json

func init() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!
	var err error
	configData, err := ioutil.ReadFile("config.yml")
	check(err)
	err = yaml.Unmarshal(configData, &config)
	check(err)
	methodChooser, err = wr.NewChooser(
		wr.Choice{Item: "netbanking", Weight: config.Netbanking.Weight},
		wr.Choice{Item: "wallet", Weight: config.Wallet.Weight},
		wr.Choice{Item: "upi", Weight: config.Upi.Weight},
		wr.Choice{Item: "card", Weight: config.Card.Weight},
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	files, err := ioutil.ReadDir("templates")
	if err != nil {
		log.Fatal(err)
	}
	templates = make(map[string]json.Json)
	for _, f := range files {
		template, _ := ioutil.ReadFile("templates/" + f.Name())
		templates[f.Name()] = json.Json{Data: template}
	}
}

func generatePayload(createdAt time.Time) []byte {
	method := methodChooser.Pick().(string)
	template := templates[fmt.Sprintf("%s_%s.json", *event, method)]
	if method == "netbanking" {
		for _, field := range config.Netbanking.Fields {
			value := field.Values[rand.Intn(len(field.Values))]
			template.Set(field.Path, []byte(value))
		}
	}
	if method == "card" {
		for _, field := range config.Card.Fields {
			value := field.Values[rand.Intn(len(field.Values))]
			template.Set(field.Path, []byte(value))
		}
	}
	if method == "upi" {
		for _, field := range config.Upi.Fields {
			value := field.Values[rand.Intn(len(field.Values))]
			template.Set(field.Path, []byte(value))
		}
	}
	if method == "wallet" {
		for _, field := range config.Wallet.Fields {
			value := field.Values[rand.Intn(len(field.Values))]
			template.Set(field.Path, []byte(value))
		}
	}
	for _, field := range config.RangeFields {
		value := rand.Intn(field.Max-field.Min) + field.Min
		template.Set(field.Path, []byte(fmt.Sprintf(`%d`, value)))

	}
	return template.Data
}
