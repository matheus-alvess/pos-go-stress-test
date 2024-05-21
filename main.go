package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(url string, requests <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }, // Desativa o acompanhamento de redirecionamentos
	}
	for range requests {
		resp, err := client.Get(url)
		if err != nil {
			results <- -1
			continue
		}
		results <- resp.StatusCode
		resp.Body.Close()
	}
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 1, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("URL é obrigatória")
		flag.Usage()
		return
	}

	requests := make(chan int, *totalRequests)
	results := make(chan int, *totalRequests)

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(*url, requests, results, &wg)
	}

	go func() {
		for i := 0; i < *totalRequests; i++ {
			requests <- i
		}
		close(requests)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	statusCounts := make(map[int]int)
	for status := range results {
		statusCounts[status]++
	}

	duration := time.Since(start)

	fmt.Printf("Tempo total gasto: %v\n", duration)
	fmt.Printf("Quantidade total de requests: %d\n", *totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", statusCounts[200])

	fmt.Println("Distribuição de códigos de status HTTP:")
	for status, count := range statusCounts {
		if status != 200 {
			if status == -1 {
				fmt.Printf("Erro de rede ou timeout: %d requests\n", count)
			} else {
				fmt.Printf("Status %d: %d requests\n", status, count)
			}
		}
	}
}
