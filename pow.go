package main

import (
	"fmt"
	"strconv"
)

// simpleHash implements the same hash function as the JavaScript version
// This is compatible with Netscape 4+ browsers
func simpleHash(str string) uint32 {
	var hash int32 = 0

	if len(str) == 0 {
		return uint32(hash)
	}

	for i := 0; i < len(str); i++ {
		char := int32(str[i])
		hash = ((hash << 5) - hash) + char
		hash = hash & hash // Convert to 32bit integer (same as JavaScript)
	}

	// Make sure we get a positive number and add some variation
	// Use absolute value like JavaScript Math.abs()
	if hash < 0 {
		hash = -hash
	}
	if hash == 0 {
		hash = 1
	}

	return uint32(hash)
}

// toHex converts a number to hex string with padding to 8 characters
func toHex(num uint32) string {
	hex := fmt.Sprintf("%x", num)
	// Pad with zeros to ensure consistent length of 8 characters
	for len(hex) < 8 {
		hex = "0" + hex
	}
	return hex
}

// hasTrailingZeros checks if the hash has the required number of trailing zeros
func hasTrailingZeros(hashNum uint32, zeros int) bool {
	hexHash := toHex(hashNum)
	trailingZeros := 0

	for i := len(hexHash) - 1; i >= 0 && trailingZeros < zeros; i-- {
		if hexHash[i] == '0' {
			trailingZeros++
		} else {
			break
		}
	}

	return trailingZeros >= zeros
}

// VerifyProofOfWork verifies that the given solution solves the proof of work challenge
// Parameters:
//   - challenge: the original challenge string
//   - solution: the solution number found by the client
//   - difficulty: the number of trailing zeros required (default: 4)
//
// Returns:
//   - bool: true if the proof of work is valid, false otherwise
func VerifyProofOfWork(challenge string, solution int, difficulty int) bool {
	if difficulty <= 0 {
		difficulty = 4 // Default difficulty
	}

	// Create the test string by concatenating challenge and solution
	testString := challenge + strconv.Itoa(solution)

	// Calculate the hash
	hash := simpleHash(testString)

	// Check if the hash has the required number of trailing zeros
	return hasTrailingZeros(hash, difficulty)
}
