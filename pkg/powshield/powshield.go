package powshield

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"math/rand"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type PoWConfig struct {
	Difficulty int
}

// GenerateNonce creates a random nonce for the PoW challenge
func GenerateNonce() string {
	return strconv.Itoa(rng.Intn(1000000))
}

// performPoW checks if the solution satisfies the PoW difficulty
func performPoW(nonce, solution string, difficulty int) bool {
	data := nonce + solution
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashStr, prefix)
}

// forwardConnection forwards data between two TCP connections (client <-> Fiber)
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
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read solution:", err)
		return
	}
	clientSolution := strings.TrimSpace(string(buffer[:n]))

	if performPoW(nonce, clientSolution, config.Difficulty) {
		// Notify client that they passed the challenge
		_, _ = conn.Write([]byte("OK\n"))
		fmt.Println("PoW validated. Forwarding connection to Fiber app")

		// Connect to Fiber server
		fiberConn, err := net.Dial("tcp", fiberAddress)
		if err != nil {
			fmt.Println("Failed to connect to Fiber server:", err)
			return
		}

		// Forward data between client and Fiber server
		forwardConnection(conn, fiberConn)
	} else {
		_, _ = conn.Write([]byte("Invalid PoW solution\n"))
	}
}
