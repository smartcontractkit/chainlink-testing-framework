package rand

import (
	"crypto/rand"
	"math/big"
)

// Int generates a random 63-bit integer that is equivalent to the math/rand.Int() using crypto/rand
// This fixes issues with gosetc seccurity lint issues
func Int() (int, error) {
	// Generate a 64-bit (8-byte) random number
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		return 0, err
	}

	// Convert bytes to a big integer
	num := new(big.Int).SetBytes(buf)

	// Clear the most significant bit to ensure the number is within 63-bit range
	// This is achieved by ANDing the number with the max 63-bit value
	max63BitValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 63), big.NewInt(1))
	num.And(num, max63BitValue)

	return int(num.Int64()), nil
}
