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
	jobs := make(chan struct{})
	results := make(chan result)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				worker(client, url, results)
			}
		}()
	}

	go func() {
		for i := 0; i < totalRequests; i++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalDuration time.Duration
	var totalOk int64
	statusCode := make(map[int]int)

	for r := range results {
		totalDuration += r.duration
		if r.err != nil {
			fmt.Printf("Erro na requisição: %v\n", r.err)
			statusCode[0]++
			continue
		}

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
