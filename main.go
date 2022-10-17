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

	r := gin.Default()

	r.Use(cors.Default())
	r.GET("/trigger", router.Trigger)

	r.Run(":8080")
}
