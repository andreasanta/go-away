package http

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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
