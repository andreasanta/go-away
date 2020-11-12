package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetReload(c *gin.Context) {

	log.Println("Inside Get Reload")

	totRecords, totTime := LoadDatabase()

	c.JSON(200, gin.H{
		"status":  "ok",
		"records": totRecords,
		"time":    totTime,
	})
}
