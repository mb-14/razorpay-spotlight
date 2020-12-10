package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/webhook", webhookEventHandler)

	r.Run()
}
