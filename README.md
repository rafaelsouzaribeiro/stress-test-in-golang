# Stress Test em Golang (MBA desafio técnico)

Sistema de CLI (Command Line Interface) desenvolvido em Go para realizar testes de carga em serviços web. A aplicação dispara múltiplas requisições HTTP concorrentes contra uma URL informada e gera um relatório detalhado ao final da execução.

## Funcionalidades

- Envio de requisições HTTP em lotes concorrentes controlados
- Relatório detalhado com tempo total, status HTTP e distribuição de erros
- Execução via linha de comando ou container Docker

## Parâmetros de Entrada

A aplicação aceita os seguintes parâmetros via linha de comando:

| Parâmetro       | Descrição                                         | Obrigatório | Exemplo                  |
|-----------------|----------------------------------------------------|-------------|---------------------------|
| `--url`         | URL do serviço a ser testado                        | Sim         | `https://reqbin.com/`      |
| `--requests`    | Número total de requisições a serem realizadas       | Sim         | `1000`                    |
| `--concurrency` | Número de chamadas simultâneas                      | Sim         | `10`                      |

## Relatório Gerado

Ao final da execução, o sistema apresenta no console:

- **Tempo total gasto** na execução
- **Quantidade total de requests** realizados
- **Quantidade de requests com status HTTP 200**
- **Quantidade de requests com status HTTP** (ex: 404, 500, etc.)

## Requisitos

- [Go](https://go.dev/dl/) 1.22+ (para execução local)
- [Docker](https://www.docker.com/) (para execução via container)

## Executando Localmente (sem Docker)

Clone o repositório e execute:

```bash
go run cmd/main.go --url=https://reqbin.com/ --requests=1000 --concurrency=10
```

## Executando via Docker (Obrigatório)

### 1. Build da imagem

```bash
docker build -t stress-test .
```

### 2. Executar o teste de carga

```bash
docker run stress-test --url=https://reqbin.com/ --requests=1000 --concurrency=10
```

## Exemplo de Saída

```
Iniciando teste de carga...
=== RELATÓRIO DE TESTE DE CARGA ===

Tempo total gasto na execução: 8.432s
Quantidade total de requests realizados: 1000
Quantidade de requests com status HTTP 200: 987
Quantidade de requests com status HTTP: 404 Quantidade: 2
```

## Arquitetura

O projeto utiliza o padrão de **lotes concorrentes** (batches), onde:

1. Um channel com buffer de tamanho `concurrency` atua como semáforo, limitando quantas requisições rodam simultaneamente
2. Um lote de requisições é disparado e o sistema aguarda **todas** finalizarem (`sync.WaitGroup`) antes de iniciar o próximo lote
3. Ao final, os resultados são agregados e o relatório é impresso no console

## Estrutura do Projeto

```
stress-test-in-golang/
├── cmd/
│   └── main.go          # Ponto de entrada da aplicação (CLI)
├── pkg/
│   └── runLoadTest/     # Lógica de execução do teste de carga
├── Dockerfile            # Build multi-stage (imagem final baseada em scratch)
├── go.mod
├── go.sum
└── README.md
```

## Tecnologias

- **Go** — linguagem principal
- **net/http** — cliente HTTP para realizar as requisições
- **sync** — controle de concorrência (WaitGroup e channels como semáforo)
- **Docker** — empacotamento multi-stage com imagem final `scratch` (mínima e sem dependências)