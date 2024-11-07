package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"ilia.com/word-of-wisdom/pkg/client"
)

func makeRequest(serverAddr, method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, "http://"+serverAddr+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.NewClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	return resp, nil
}

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		panic("SERVER_ADDR is not set")
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := makeRequest(serverAddr, "GET", "/quote")
			if err != nil {
				fmt.Println("Request error:", err)
				return
			}
			defer resp.Body.Close()

			respRaw, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Failed to read response body:", err)
				return
			}

			fmt.Printf("Response Status: %s , Body: %s\n", resp.Status, string(respRaw))
		}()
	}

	wg.Wait()
}
