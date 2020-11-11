package http

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	status    string
	reason    string `json:",omitempty"`
	timestamp uint   `json:",omitempty"`
}

func GetIp(c *gin.Context) {
	c.JSON(200, Response{
		status: "ok",
	})
}
