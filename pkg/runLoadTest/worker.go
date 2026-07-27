package runloadtest

import (
	"net/http"
	"time"
)

func worker(client *http.Client, url string, results chan<- result) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		results <- result{statusCode: 0, duration: time.Since(start), err: err}
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36")

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		results <- result{statusCode: 0, duration: duration, err: err}
		return
	}

	defer resp.Body.Close()
	results <- result{statusCode: resp.StatusCode, duration: duration, err: nil}
}
