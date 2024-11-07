package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ilia.com/word-of-wisdom/pkg/powshield"
)

var client = NewClient()

// powHandshake performs the PoW handshake, establishes a validated TCP connection.
func powHandshake(conn net.Conn) (net.Conn, error) {
	reader := bufio.NewReader(conn)
	challenge, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge: %w", err)
	}

	parts := strings.Fields(challenge)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid challenge format")
	}
	nonce := strings.Split(parts[1], ":")[1]
	difficulty, _ := strconv.Atoi(strings.Split(parts[2], ":")[1])

	solution := powshield.Solve(nonce, difficulty)
	fmt.Printf("Solved PoW: %s\n", solution)

	_, err = conn.Write([]byte(solution + "\n"))
	if err != nil {
		return nil, fmt.Errorf("failed to send solution: %w", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(response) != "OK" {
		return nil, fmt.Errorf("PoW validation failed or server closed connection")
	}

	return conn, nil
}

func configureAPIConnection(conn *net.TCPConn) {
	conn.SetNoDelay(true)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
}

func GetClient() *http.Client {
	return client
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{

			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				fmt.Printf("Dialing to %s\n", addr)
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
}
