package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"ilia.com/word-of-wisdom/pkg/powshield"
)

var quotes = []string{
	"The only way to do great work is to love what you do. - Steve Jobs",
	"Life is what happens when you're busy making other plans. - John Lennon",
	"Success is not final, failure is not fatal. - Winston Churchill",
	"Be the change you wish to see in the world. - Mahatma Gandhi",
	"Stay hungry, stay foolish. - Steve Jobs",
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomQuote() string {
	return quotes[rng.Intn(len(quotes))]
}

func main() {

	// Initialize Fiber app
	app := fiber.New()

	// Handle all incoming requests to the root
	app.Get("/quote", func(c *fiber.Ctx) error {
		return c.SendString(GetRandomQuote())
	})
	// PoW difficulty level
	config := powshield.PoWConfig{Difficulty: 4}

	// Fiber app address to forward connections to
	fiberAddress := "localhost:3000" // Fiber server running on port 3000

	// Start TCP server for PoW validation
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Failed to start TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server started on :8081")

	go func() {
		for {
			// Accept incoming TCP connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Failed to accept connection:", err)
				continue
			}

			// Handle the connection with PoW and forwarding
			go powshield.HandleTCPConnection(conn, config, fiberAddress)
		}
	}()

	err = app.Listen(":3000")
	if err != nil {
		fmt.Println("Failed to start Fiber server:", err)
	}
}
