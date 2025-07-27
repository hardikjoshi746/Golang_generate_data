# Golang Generate Data API – Manifold Labs Takehome Project

## Overview

This project implements a high-performance REST API in Go that:
- Accepts POST requests from users.
- Simulates a random delay (100ms to 50,000ms).
- Generates word data proportional to the delay.
- Tracks user quotas in a MySQL database.
- Supports 5000+ concurrent requests reliably.

## Features

- `POST /generate-data`  
  Generates and returns user-specific data based on a simulated delay. Tracks usage quotas.

- `GET /user/requests`  
  Returns all requests and data generated for a specific user.

- `GET /user/stats`  
  Returns user statistics like total requests, average delay, words used, and words left.

## Tech Stack

- **Go (Gin)** – Lightweight, high-performance web framework.
- **MySQL** – Persistent storage for user quotas and request history.
- **Redis** *(Optional for performance optimization)* – Can be added for caching responses or quota tracking.
- **Docker** – Containerized for portability and reproducibility.
- **Postman** – Used for manual testing.

## Database Schema

### `users`
| user_id | word_used | word_left | created_at |
|---------|-----------|-----------|------------|

### `requests`
| id (PK) | user_id (FK) | data (TEXT) | duration (ms) | created_at |

## Business Logic

- Each user starts with **1,000,000 word quota**.
- Each request consumes `delay(ms) / 6` words.
- If the user exceeds their quota, the request is **rejected**.
- All word usage is transaction-safe via `FOR UPDATE` row locks to prevent race conditions.

## Quantitative Outcomes

- ✅ **Load Test @ 2000 Requests, 500 Concurrency**  
Success: 2000
Failures: 0
Total time: ~2m4s
Avg per req: ~66ms


- ✅ **Load Test @ 5000 Requests, 1000 Concurrency**  
Success: 5000
Failures: 0
Total time: ~2m39s
Avg per req: ~31.94ms


- ⚙️ Scales to **5000+ concurrent** users with proper DB connection tuning.

## Running the App

### 1. Start MySQL locally or via Docker
Make sure you have a running MySQL instance.

### 2. Run the application

```bash
go run main.go
```
3. Run the load test (optional)
```bash
go run load_test_main.go
```


👤 Hardik Joshi
