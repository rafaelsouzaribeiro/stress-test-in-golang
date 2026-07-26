package main

import (
	"flag"
	"fmt"
	"os"

	runloadtest "github.com/rafaelsouzaribeiro/stress-test-in-golang/pkg/runLoadTest"
)

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

	runloadtest.RunLoadTest(*url, *requests, *concurrency)

}
