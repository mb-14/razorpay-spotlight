package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mb-14/rzp-spotlight/webhook/json"
	"github.com/mb-14/rzp-spotlight/webhook/rzp"
)

const (
	DBName = "rzpftx"
)

func webhookEventHandler(c *gin.Context) {
	client := influxdb2.NewClient("http://localhost:8086", "")
	writeAPI := client.WriteAPIBlocking("", DBName)
	defer client.Close()
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json := json.Json{Data: payload}
	p := rzp.ProcessPayloadJson(json)
	err = writeAPI.WritePoint(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
