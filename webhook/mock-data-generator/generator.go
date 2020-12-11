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
		wr.Choice{Item: "netbanking", Weight: config.MethodDistribution.Netbanking},
		wr.Choice{Item: "wallet", Weight: config.MethodDistribution.Wallet},
		wr.Choice{Item: "upi", Weight: config.MethodDistribution.Upi},
		wr.Choice{Item: "card", Weight: config.MethodDistribution.Card},
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
	return template.Data
}
