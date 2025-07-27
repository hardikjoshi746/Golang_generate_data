package main

import (
	"net/http"

	"example.com/rest/database"
	"example.com/rest/handlers"
	"example.com/rest/redis"
	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	redis.InitRedis()
	database.InitDB()

	server.GET("/health", healthCheck)
	server.POST("/generate-data", handlers.GenerateDataHandler)
	server.GET("/user/requests", handlers.ViewUserRequests)
	server.GET("/user/stats", handlers.UserStatsHandler)
	server.GET("/redis-test", func(c *gin.Context) {
		val, err := redis.Rdb.Get(redis.Ctx, "test-key").Result()
		if err != nil {
			// Set it if not exists
			redis.Rdb.Set(redis.Ctx, "test-key", "Hello from Redis!", 0)
			c.JSON(200, gin.H{"message": "Key not found, setting now", "set_value": "Hello from Redis!"})
			return
		}
		c.JSON(200, gin.H{"redis_value": val})
	})

	server.Run(":8080")

}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
