package main

import (
    "log"
    "retail_pulse/config"
    "retail_pulse/db"
    "retail_pulse/handlers"

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
