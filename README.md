# Golang Generate Data API 

## Overview

This high-performance REST API accepts data generation requests, simulates a variable delay, tracks word usage per user, and enforces user quotas.

Built with **Go**, **MySQL**, and **Redis**, the app supports over **5000 concurrent requests**, includes data persistence, and is fully containerized.

---

## Features

- `POST /generate-data`  
  Simulates word generation based on random delay and deducts word quota from users. Enforces rate-limiting.

- `GET /user/requests`  
  Fetches request history for a given user.

- `GET /user/stats`  
  Returns statistics like total requests, words used, and average delay.

---

## Database Schema

### `users`  
Tracks each user's word usage.

| Column     | Type         | Notes                     |
|------------|--------------|---------------------------|
| user_id    | VARCHAR      | Primary Key               |
| word_used  | INT          | Increments per request    |
| words_left | INT          | Starts at 1,000,000       |
| created_at | TIMESTAMP    | Defaults to now           |

### `requests`  
Stores request history and timing.

| Column     | Type         | Notes                     |
|------------|--------------|---------------------------|
| request_id | INT (auto)   | Primary Key               |
| user_id    | VARCHAR      | Foreign Key to `users`    |
| data       | TEXT         | Generated words           |
| duration   | INT          | Simulated delay in ms     |
| created_at | TIMESTAMP    | Defaults to now           |

---

## Business Logic

- Each user starts with **1,000,000** word quota.
- A request incurs a **random delay (100–50,000ms)**.
- Words generated = `delay / 100` (10ms = 1 word).
- If the user doesn't have enough words left, the request is rejected.
- Rate limiting of **10 requests per minute** enforced via Redis.
- All DB writes are **transactional with row locking**.

---

## Tech Stack

- **Go (Gin)** — High-performance REST framework
- **MySQL** — Persistent request and user quota store
- **Redis** — Rate limiter for request throttling
- **Docker** — Containerization for all services

---

## Folder Structure

```
.
├── main.go               # Entry point
├── handlers/             # Route handlers
├── database/             # MySQL logic
├── redis/                # Redis-based rate limiter
├── config/               # .env reader
├── Dockerfile
├── docker-compose.yml
├── .env
└── init.sql              # DB bootstrap
```

---

## Running the App (Dockerized)

### 1. Clone the Repo

```bash
git clone https://github.com/yourusername/golang-generate-data
cd golang-generate-data
```

### 2. Set Up Environment

Create a `.env` file:

```env
DB_HOST=mysql
DB_PORT=3306
DB_USER=youruser
DB_PASSWORD=yourpassword
DB_NAME=golang_db

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

SERVER_PORT=8080
ENVIRONMENT=production
```

For local development, use `.env.local` with `localhost`.

### 3. Start Everything

```bash
docker-compose up --build
```

### 4. Test Health Check (Optional)

```bash
curl http://localhost:8080/health
```

---

## 🔍 API Examples

### ▶️ Generate Data

```bash
curl -X POST http://localhost:8080/generate-data -H "X-User-id: 123"
```

### 📜 User Request History

```bash
curl -X GET http://localhost:8080/user/requests -H "X-User-id: 123"
```

### 📊 User Stats

```bash
curl -X GET http://localhost:8080/user/stats -H "X-User-id: 123"
```

---

## Performance

The API is optimized to handle high concurrency using:

- Gin’s goroutine-friendly handler model
- MySQL connection pooling
- Redis-based rate limiter
- Transaction-safe quota enforcement

### Load Test Results

| Scenario                     | Requests | Concurrency | Avg Time | Success Rate |
|-----------------------------|----------|-------------|----------|--------------|
| Moderate load               | 2000     | 500         | ~66ms    | 100%         |
| Heavy load                  | 5000     | 1000        | ~32ms    | 100%         |

---

## Manual Testing

You can also use [Postman](https://www.postman.com/) or cURL to interact with endpoints using the `X-User-id` header.

---

## Author

**Hardik Joshi**
