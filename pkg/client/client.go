package client

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"math/rand"
)

// Create a new HTTP client and use the existing TCP connection as the transport layer
var client = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			if tcpConn, ok := conn.(*net.TCPConn); ok {
				configureAPIConnection(tcpConn)
			}

			conn, err = powHandshake(conn)
			if err != nil {
				fmt.Println("PoW handshake error:", err)
				return nil, err
			}
			return conn, nil
		},
	},
}

// powHandshake performs the PoW handshake, establishes a validated TCP connection.
func powHandshake(conn net.Conn) (net.Conn, error) {
	// Read PoW challenge from the server
	reader := bufio.NewReader(conn)
	challenge, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge: %w", err)
	}

	// Parse the challenge for nonce and difficulty
	parts := strings.Fields(challenge)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid challenge format")
	}
	nonce := strings.Split(parts[1], ":")[1]
	difficulty, _ := strconv.Atoi(strings.Split(parts[2], ":")[1])

	// Solve the challenge
	solution := performPoW(nonce, difficulty)

	// Send solution to the server
	_, err = conn.Write([]byte(solution + "\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to send solution: %w", err)
	}

	// Await server validation response
	response, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(response) != "OK" {
		return nil, fmt.Errorf("PoW validation failed or server closed connection")
	}

	// Return the validated TCP connection
	fmt.Println("PoW validated, TCP connection established")
	return conn, nil
}

// performPoW solves the challenge by finding a valid nonce solution.
func performPoW(nonce string, difficulty int) string {
	prefix := strings.Repeat("0", difficulty)
	for {
		solution := strconv.Itoa(rand.Intn(1000000))
		hash := sha256.Sum256([]byte(nonce + solution))
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, prefix) {
			return solution
		}
	}
}

func configureAPIConnection(conn *net.TCPConn) {
	conn.SetNoDelay(true)                                  // Reduce latency for responses
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))  // Prevent slow clients
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second)) // Prevent slow responses
}

func GetClient() *http.Client {
	return client
}
