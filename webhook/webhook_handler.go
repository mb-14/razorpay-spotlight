package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	DBName = "rzpftx"
)

var (
	client   = influxdb2.NewClient("http://localhost:8086", "")
	writeAPI = client.WriteAPIBlocking("", DBName)
)

// Events
const (
	PaymentAuthorized = "payment.authorized"
	PaymentFailed     = "payment.failed"
)

// Methods
const (
	Netbanking = "netbanking"
	Wallet     = "wallet"
	UPI        = "upi"
	Card       = "card"
)

func webhookEventHandler(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json := Json{payload}
	var event string
	if value, err := json.GetString("event"); err == nil {
		event = value
	}
	if event == PaymentAuthorized {
		err = writePaymentEvent(c.Request.Context(), json, "payment_authorized")
	} else if event == PaymentFailed {
		err = writePaymentEvent(c.Request.Context(), json, "payment_failed")
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}

func writePaymentEvent(ctx context.Context, json Json, measurement string) error {
	amount, _ := json.GetInt("payload.payment.entity.amount")
	createdAt, _ := json.GetTime("payload.payment.entity.created_at")
	p := influxdb2.NewPoint(measurement,
		addTags(json),
		map[string]interface{}{"amount": amount},
		createdAt)
	return writeAPI.WritePoint(ctx, p)
}

func addTags(p Json) map[string]string {
	tags := make(map[string]string)
	// Common tags
	method, _ := p.GetString("payload.payment.entity.method")
	currency, _ := p.GetString("payload.payment.entity.currency")
	tags["method"] = method
	tags["currency"] = currency

	if method == Netbanking {
		tags["bank"], _ = p.GetString("payload.payment.entity.bank")
	}

	if method == Wallet {
		// TODO add wallet name
	}

	if method == UPI {
		// TODO - Exract psp from VPA and add as tag (@okhdfcbank, @oksbi, @paytm, @upi)
		vpa, _ = p.GetString("payload.payment.entity.vpa")
		vpaString := strings.Split(vpa, "@")
		tags["psp"] = vpaString[1]
	}

	if method == Card {
		// TODO Add card details
	}

	return tags
}
