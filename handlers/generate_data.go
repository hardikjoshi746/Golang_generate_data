package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardikjoshi746/Golang_generate_data/database"
	"github.com/hardikjoshi746/Golang_generate_data/redis"
)

func GenerateDataHandler(c *gin.Context) {
	userID := c.GetHeader("X-User-id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}

	limited, err := redis.RateLimit(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
		return
	}
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		return
	}

	delayMs := rand.Intn(50000) + 100 // 100 to 50000
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	// Step 3: Simulate word usage and remaining
	wordsUsed := int(delayMs / 100) // 1 words per millisecond
	// wordsLeft := 1000000 - wordsUsed

	// if wordsLeft < 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "user has no words left"})
	// 	return
	// } // if the user has no words left, return an error

	// Generate words by repeating the dummy words array as needed
	result := Result(wordsUsed)

	// Ensure user exists in database
	if err := database.EnsureUserExists(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save request to database
	// err := database.SaveRequest(userID, result, delay)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Update user words in database
	// err = database.UpdateUserWords(userID, wordsUsed)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	err = database.ProcessUserRequestTx(userID, result, wordsUsed, delayMs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())

	wordsLeft, err := database.GetUserWordBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Do something with userID (e.g., check in DB, simulate processing, etc.)
	c.JSON(http.StatusOK, gin.H{
		// "message": "POST request received",
		// "user_id":    userID,
		// "status":     "success",
		"time to respond(in ms)": delayMs,
		"data":                   result,
		"words_used":             wordsUsed,
		"words_left":             wordsLeft,
	})
}

func ViewUserRequests(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-Id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-User-Id"})
		return
	}

	requests, err := database.GetUserRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
