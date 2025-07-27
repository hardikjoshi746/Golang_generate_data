package main

import (
	"net/http"

	"example.com/rest/database"
	"example.com/rest/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	database.InitDB()

	server.GET("/health", healthCheck)
	server.POST("/generate-data", handlers.GenerateDataHandler)
	server.GET("/user/requests", handlers.ViewUserRequests)
	server.GET("/user/stats", handlers.UserStatsHandler)

	server.Run(":8080")

}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
