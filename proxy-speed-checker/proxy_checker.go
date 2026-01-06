package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type ProxyResult struct {
	Proxy              string
	AverageRequestTime float64
}

func main() {
	proxies, err := loadProxies("proxy.txt")
	if err != nil {
		fmt.Println("Error loading proxies:", err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan ProxyResult, len(proxies))

	numChecks := 5
	maxConcurrentChecks := 300
	semaphore := make(chan struct{}, maxConcurrentChecks)
	progress := make(chan int, len(proxies))

	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			semaphore <- struct{}{}        // acquire semaphore
			defer func() { <-semaphore }() // release semaphore

			requestTime := testProxy(proxy, numChecks)
			results <- ProxyResult{Proxy: proxy, AverageRequestTime: requestTime}
			progress <- 1
		}(proxy)
	}

	go func() {
		wg.Wait()
		close(results)
		close(progress)
	}()

	go trackProgress(progress, len(proxies))

	if err := saveResults("results.csv", results); err != nil {
		fmt.Println("Error saving results:", err)
	}
}

func loadProxies(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	return proxies, scanner.Err()
}

func testProxy(proxy string, numChecks int) float64 {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return 0
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	var totalRequestTime float64

	for i := 0; i < numChecks; i++ {
		startRequest := time.Now()
		resp, err := httpClient.Get("https://github.com")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
		}
		totalRequestTime += time.Since(startRequest).Seconds()
		time.Sleep(400 * time.Millisecond)
	}

	return totalRequestTime / float64(numChecks)
}

func saveResults(filename string, results chan ProxyResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Proxy", "AverageRequestTime"})

	for result := range results {
		writer.Write([]string{
			result.Proxy,
			fmt.Sprintf("%.2f", result.AverageRequestTime),
		})
	}

	return nil
}

func trackProgress(progress chan int, total int) {
	completed := 0
	for range progress {
		completed++
		fmt.Printf("Progress: %d/%d\n", completed, total)
	}
}
