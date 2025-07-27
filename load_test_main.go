package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	targetURL   = "http://localhost:8080/generate-data"
	concurrency = 1000
	totalReqs   = 5000
	userID      = "testuser-load"
)

func load_test_main() {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	start := time.Now()
	var success, failure int
	var lock sync.Mutex

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			req, _ := http.NewRequest("POST", targetURL, bytes.NewBuffer([]byte{}))
			req.Header.Set("X-User-Id", fmt.Sprintf("%d", 1000+i))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				lock.Lock()
				failure++
				lock.Unlock()
				fmt.Printf("❌ Request %d failed: %v\n", i, err)
				return
			}
			resp.Body.Close()

			if resp.StatusCode == 200 {
				lock.Lock()
				success++
				lock.Unlock()
			} else {
				lock.Lock()
				failure++
				lock.Unlock()
				fmt.Printf("❌ Request %d failed with status: %d\n", i, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("========== Load Test Summary ==========")
	fmt.Printf("Total Requests: %d\n", totalReqs)
	fmt.Printf("Concurrency:    %d\n", concurrency)
	fmt.Printf("Success:        %d\n", success)
	fmt.Printf("Failures:       %d\n", failure)
	fmt.Printf("Total time:     %v\n", elapsed)
	fmt.Printf("Avg per req:    %.2fms\n", float64(elapsed.Milliseconds())/float64(totalReqs))
}

func main() {
	load_test_main()
}
