package http

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Response struct {
	Status    string `json:"status"`
	Reason    string `json:"description,omitempty"`
	Timestamp uint   `json:"timestamp,omitempty"`
}

func GetIp(c *gin.Context) {
	
	log.Println("Inside Get IP")

	c.JSON(200, Response{
		Status: "ok",
	})
}
