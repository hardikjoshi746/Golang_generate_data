package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hardikjoshi746/Golang_generate_data/redis"
)

type UserRequest struct {
	UserID   int
	Data     string
	Duration int
}

// UpdateUserWords updates the user's word usage
func UpdateUserWords(userID string, wordsUsed int) error {
	query := `
		UPDATE users
		SET word_used = word_used + ?, words_left = words_left - ?
		WHERE user_id = ?
	`
	_, err := DB.Exec(query, wordsUsed, wordsUsed, userID)
	if err != nil {
		return fmt.Errorf("failed to update user words: %v", err)
	}
	return nil
}

// SaveRequest saves the generated data to the requests table
func SaveRequest(userID string, data string, duration int) error {
	query := `
		INSERT INTO requests (user_id, data, duration)
		VALUES (?, ?, ?)
	`
	_, err := DB.Exec(query, userID, data, duration)
	if err != nil {
		return fmt.Errorf("failed to save request: %v", err)
	}
	return nil
}

func EnsureUserExists(userID string) error {
	query := `
		INSERT INTO users (user_id, word_used, words_left)
		VALUES (?, 0, 1000000)
		ON DUPLICATE KEY UPDATE user_id = user_id
	`
	_, err := DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to ensure user exists: %v", err)
	}
	return nil
}

func ProcessUserRequestTx(userID string, result string, wordsUsed int, duration int) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Step 1: Lock user row
	var wordLeft int
	query := `SELECT words_left FROM users WHERE user_id = ? FOR UPDATE`
	err = tx.QueryRow(query, userID).Scan(&wordLeft)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to lock user row: %v", err)
	}

	// Step 2: Check if user has enough words
	if wordLeft < wordsUsed {
		tx.Rollback()
		return fmt.Errorf("user has insufficient words: has %d, needs %d", wordLeft, wordsUsed)
	}

	// Step 3: Update user word usage
	updateUser := `
		UPDATE users
		SET word_used = word_used + ?, words_left = words_left - ?
		WHERE user_id = ?
	`
	_, err = tx.Exec(updateUser, wordsUsed, wordsUsed, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user words: %v", err)
	}

	// Step 4: Save request
	insertRequest := `
		INSERT INTO requests (user_id, data, duration)
		VALUES (?, ?, ?)
	`
	_, err = tx.Exec(insertRequest, userID, result, duration)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert request: %v", err)
	}

	// Step 5: Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func GetUserWordBalance(userID string) (int, error) {
	// Check Redis cache first
	cacheKey := "words_left:" + userID
	val, err := redis.Rdb.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		fmt.Printf("✅ Cache HIT for key: %s, value: %s\n", cacheKey, val)
		wordsLeft, _ := strconv.Atoi(val)
		return wordsLeft, nil
	} else {
		fmt.Printf("❌ Cache MISS for key: %s, reason: %v\n", cacheKey, err)
	}

	// Fallback to DB
	var wordsLeft int
	err = DB.QueryRow("SELECT words_left FROM users WHERE user_id = ?", userID).Scan(&wordsLeft)
	if err != nil {
		return 0, err
	}

	// Set cache
	redis.Rdb.Set(redis.Ctx, cacheKey, wordsLeft, time.Minute*5)
	return wordsLeft, nil
}

func GetUserRequests(userID int) ([]UserRequest, error) {
	query := `SELECT user_id, data, duration FROM requests WHERE user_id = ?`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user requests: %v", err)
	}
	defer rows.Close()

	var results []UserRequest
	for rows.Next() {
		var ur UserRequest
		if err := rows.Scan(&ur.UserID, &ur.Data, &ur.Duration); err != nil {
			return nil, err
		}
		results = append(results, ur)
	}
	return results, nil
}

func GetUserStats(userID string) (map[string]interface{}, error) {
	query := `
		SELECT u.user_id, u.word_used, u.words_left,
		       COUNT(r.id) as total_requests,
		       IFNULL(AVG(r.duration), 0) as avg_duration
		FROM users u
		LEFT JOIN requests r ON u.user_id = r.user_id
		WHERE u.user_id = ?
		GROUP BY u.user_id, u.word_used, u.words_left
		`
	row := DB.QueryRow(query, userID)

	var uid string
	var used, left, total int
	var avg float64

	err := row.Scan(&uid, &used, &left, &total, &avg)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %v", err)
	}
	return map[string]interface{}{
		"user_id":        uid,
		"words_used":     used,
		"words_left":     left,
		"total_requests": total,
		"avg_delay_ms":   avg,
	}, nil
}
