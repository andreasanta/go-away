package http

import (
	"os"

	"github.com/gin-gonic/gin"
)

func RunServer() {

	loadConfig()

	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		v1.GET("/ip/:ip", GetIp)
		//v1.GET("/reload", GetReload)
	}

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
