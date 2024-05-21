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
		Timeout: 1 * time.Second,
	}
	for range requests {
		resp, err := client.Get(url)
		if err != nil {
			//fmt.Println("Erro ao fazer a requisição:", err)
			continue
		}
		if resp.StatusCode != 0 {
			results <- resp.StatusCode
		}
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
	errorStatusCodes := make(map[int]bool)
	for status := range results {
		statusCounts[status]++
		if status != 200 {
			errorStatusCodes[status] = true
		}
	}

	duration := time.Since(start)

	fmt.Printf("Tempo total gasto: %v\n", duration)
	fmt.Printf("Quantidade total de requests: %d\n", *totalRequests)
	for status, count := range statusCounts {
		if status == 0 {
			fmt.Printf("Erro ao fazer request: %d\n", count)
		} else {
			fmt.Printf("Status %d: %d requests\n", status, count)
		}
	}

	if len(errorStatusCodes) > 0 {
		fmt.Println("Códigos de status HTTP de erro capturados:")
		for status := range errorStatusCodes {
			fmt.Printf("Status %d\n", status)
		}
	} else {
		fmt.Println("Nenhum código de status HTTP de erro encontrado.")
	}
}
