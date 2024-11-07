package powshield

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// old choice

// // GenerateNonce creates a random nonce for the PoW challenge
// func GenerateNonce() string {
// 	return strconv.Itoa(rng.Intn(1000000))
// }

// // CheckSolution checks if the solution satisfies the PoW difficulty
// func CheckSolution(nonce, solution string, difficulty int) bool {
// 	data := nonce + solution
// 	hash := sha256.Sum256([]byte(data))
// 	hashStr := hex.EncodeToString(hash[:])
// 	prefix := strings.Repeat("0", difficulty)
// 	return strings.HasPrefix(hashStr, prefix)
// }

// // Solve solves the challenge by finding a valid nonce solution.
// func Solve(nonce string, difficulty int) string {
// 	prefix := strings.Repeat("0", difficulty)
// 	for {
// 		solution := strconv.Itoa(rand.Intn(1000000))
// 		hash := sha256.Sum256([]byte(nonce + solution))
// 		hashStr := hex.EncodeToString(hash[:])
// 		if strings.HasPrefix(hashStr, prefix) {
// 			return solution
// 		}
// 	}
// }

// Algo https://ru.wikipedia.org/wiki/Scrypt

// GenerateNonce creates a random nonce for the PoW challenge.
func GenerateNonce() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return num.String()
}

// CheckSolution verifies if the solution satisfies the PoW difficulty using Scrypt.
func CheckSolution(nonce, solution string, difficulty int) bool {
	data := nonce + solution
	scryptHash, err := scrypt.Key([]byte(data), []byte("somesalt"), 1<<8, 2, 1, 32)
	if err != nil {
		return false
	}
	hash := sha256.Sum256(scryptHash)
	hashStr := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashStr, prefix)
}

// Solve attempts to find a solution to the PoW challenge by brute-forcing solutions.
func Solve(nonce string, difficulty int) string {
	prefix := strings.Repeat("0", difficulty)

	for i := 0; i < 1000000; i++ {
		solution := fmt.Sprintf("%d", i)
		data := nonce + solution

		// Generate a memory-hard hash using Scrypt
		scryptHash, err := scrypt.Key([]byte(data), []byte("somesalt"), 1<<8, 2, 1, 32)
		if err != nil {
			fmt.Println("Error generating Scrypt hash:", err)
			continue
		}

		hash := sha256.Sum256(scryptHash)
		hashStr := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashStr, prefix) {
			return solution
		}
	}
	return ""
}
