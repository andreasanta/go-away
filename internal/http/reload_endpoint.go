package http

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GetReload(c *gin.Context) {

	log.Println("Inside Get Reload")

	LoadDatabase()

	c.JSON(200, Response{
		Status: "ok",
	})
}
