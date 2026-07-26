package runloadtest

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type result struct {
	statusCode int
	duration   time.Duration
	err        error
}

func RunLoadTest(url string, totalRequests, concurrency int) {
	requestsPerWorker := totalRequests / concurrency
	remainder := totalRequests % concurrency

	results := make(chan result, totalRequests)
	client := &http.Client{}
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		for j := 0; j < requestsPerWorker; j++ {
			wg.Add(1)
			go worker(client, url, results, &wg)
		}
	}

	for i := 0; i < remainder; i++ {
		wg.Add(1)
		go worker(client, url, results, &wg)
	}

	wg.Wait()
	close(results)

	var totalDuration time.Duration
	var totalOk int64 = 0
	var statusCode = make(map[int]int)
	for r := range results {
		totalDuration += r.duration
		if r.statusCode == http.StatusOK {
			totalOk++
			continue
		}
		statusCode[r.statusCode]++
	}

	fmt.Println("=== RELATÓRIO DE TESTE DE CARGA ===")
	fmt.Println()
	fmt.Printf("Tempo total gasto na execução: %v\n", totalDuration)
	fmt.Printf("Quantidade total de requests realizados: %v\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200. %v\n", totalOk)

	for code, count := range statusCode {
		fmt.Printf("Quantidade de requests com status HTTP: %d Quantidade: %v\n", code, count)
	}
}
