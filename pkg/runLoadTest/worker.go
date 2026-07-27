package runloadtest

import (
	"net/http"
	"time"
)

func worker(client *http.Client, url string, results chan<- result) {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		results <- result{statusCode: 0, duration: duration, err: err}
		return
	}

	resp.Body.Close()
	results <- result{statusCode: resp.StatusCode, duration: duration, err: nil}
}
