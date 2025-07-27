package handlers

import (
	"net/http"

	"example.com/rest/database"
	"github.com/gin-gonic/gin"
)

func UserStatsHandler(c *gin.Context) {
	userId := c.GetHeader("X-User-Id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id not found"})
		return
	}
	stats, err := database.GetUserStats(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
