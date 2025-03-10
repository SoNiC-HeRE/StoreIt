package main

import (
	"log"
	"storeit/config"
	"storeit/database"
	"storeit/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.ConnectMongoDB()

	router := gin.Default()
	router.POST("/api/submit/", handlers.SubmitJob)
	router.GET("/api/status", handlers.GetJobStatus)

	log.Println("Server running on :8080")
	router.Run(":8080")
}
