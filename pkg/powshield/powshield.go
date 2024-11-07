package powshield

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"math/rand"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type PoWConfig struct {
	Difficulty int
}

// forwardConnection forwards data between two TCP connections (client <-> Fiberapp)
func forwardConnection(clientConn, fiberConn net.Conn) {
	defer clientConn.Close()
	defer fiberConn.Close()

	// Relay data between client and Fiber in both directions
	go io.Copy(clientConn, fiberConn)
	io.Copy(fiberConn, clientConn)
}

// HandleTCPConnection manages PoW and forwards connection if validated
func HandleTCPConnection(conn net.Conn, config PoWConfig, fiberAddress string) {
	defer conn.Close()

	// Step 1: Generate and send the PoW challenge
	nonce := GenerateNonce()
	challenge := fmt.Sprintf("SYN-ACK nonce:%s difficulty:%d\n", nonce, config.Difficulty)
	_, err := conn.Write([]byte(challenge))
	if err != nil {
		fmt.Println("Failed to send challenge:", err)
		return
	}

	// Step 2: Read and verify the solution
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to read solution: %s , buffer: %s\n", err, string(buffer[:n]))
		return
	}
	clientSolution := strings.TrimSpace(string(buffer[:n]))

	if CheckSolution(nonce, clientSolution, config.Difficulty) {
		_, _ = conn.Write([]byte("OK\n"))
		fmt.Println("PoW validated. Forwarding connection to Fiber app")

		fiberConn, err := net.Dial("tcp", fiberAddress)
		if err != nil {
			fmt.Println("Failed to connect to Fiber server:", err)
			return
		}

		forwardConnection(conn, fiberConn)
	} else {
		_, _ = conn.Write([]byte("Invalid PoW solution\n"))
	}
}
