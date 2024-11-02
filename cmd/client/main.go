package main

import (
	"fmt"
	"io"
	"net/http"

	"ilia.com/word-of-wisdom/pkg/client"
)

// Config for connecting to the server
const (
	serverAddr = "localhost:8081" // PoW server with TCP forwarding to Fiber
	difficulty = 4                // Expected PoW difficulty from the server
)

// makeRequest sends an HTTP request over a validated TCP connection.
func makeRequest(method, path string) (*http.Response, error) {

	// Construct the HTTP request
	req, err := http.NewRequest(method, "http://localhost:8081"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Perform the HTTP request
	resp, err := client.GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	return resp, nil
}

func main() {
	// Establish the PoW handshake and validated connection

	// Example HTTP GET request
	resp, err := makeRequest("GET", "/quote")
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

	// Print response status
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(respRaw))
}
