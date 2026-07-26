package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type result struct {
	statusCode int
	duration   time.Duration
	err        error
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requisições")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("Erro: --url é obrigatório")
		os.Exit(1)
	}
	if *requests <= 0 {
		fmt.Println("Erro: --requests deve ser maior que zero")
		os.Exit(1)
	}
	if *concurrency <= 0 {
		fmt.Println("Erro: --concurrency deve ser maior que zero")
		os.Exit(1)
	}
	if *concurrency > *requests {
		*concurrency = *requests
	}

	fmt.Printf("Iniciando teste de carga...\n")

	runLoadTest(*url, *requests, *concurrency)

}

func runLoadTest(url string, totalRequests, concurrency int) {
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
		fmt.Printf("Quantidade de requests com status HTTP %d: %v\n", code, count)
	}
}

func worker(client *http.Client, url string, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
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
