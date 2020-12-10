package main

import (
	"net/http"

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
		amount, _ := json.GetInt("payload.payment.entity.amount")
		method, _ := json.GetString("payload.payment.entity.method")
		createdAt, _ := json.GetTime("payload.payment.entity.created_at")
		p := influxdb2.NewPoint(PaymentAuthorized,
			map[string]string{"method": method},
			map[string]interface{}{"amount": amount},
			createdAt)
		err = writeAPI.WritePoint(c.Request.Context(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
