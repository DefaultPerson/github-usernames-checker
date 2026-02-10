package utils

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// CheckProxies tests each proxy with a single request to github.com
// and returns only the working ones.
func CheckProxies(proxies []string, timeout time.Duration) []string {
	if len(proxies) == 0 {
		return proxies
	}

	log.Info().Msgf("Checking %d proxies...", len(proxies))

	type result struct {
		proxy string
		alive bool
	}

	results := make(chan result, len(proxies))
	semaphore := make(chan struct{}, MaxGoroutines)
	var wg sync.WaitGroup
	var checked atomic.Int32

	for _, p := range proxies {
		wg.Add(1)
		go func(raw string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			alive := testProxyAlive(raw, timeout)
			results <- result{proxy: raw, alive: alive}

			done := checked.Add(1)
			if done%10 == 0 || int(done) == len(proxies) {
				log.Info().Msgf("Proxy check progress: %d/%d", done, len(proxies))
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var alive []string
	for r := range results {
		if r.alive {
			alive = append(alive, r.proxy)
		}
	}

	log.Info().Msgf("Proxy check done: %d/%d alive", len(alive), len(proxies))
	return alive
}

func testProxyAlive(raw string, timeout time.Duration) bool {
	normalized := NormalizeProxy(raw)

	var proxyURL *url.URL
	var err error

	if strings.HasPrefix(normalized, "socks5://") {
		proxyURL, err = url.Parse(normalized)
	} else {
		proxyURL, err = url.Parse("http://" + normalized)
	}
	if err != nil {
		return false
	}

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   timeout,
	}

	resp, err := client.Get("https://github.com")
	if err != nil {
		return false
	}
	resp.Body.Close()

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound
}
