package main

import (
	"cherry-ipfs-client/router"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	s3_client := router.Default()

	r := gin.Default()

	r.Use(cors.Default())
	r.GET("/trigger", s3_client.Trigger)

	r.Run(":8081")
}
