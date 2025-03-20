package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58"
	"strconv"
	"time"
)

const lenOfPart = 7

// GenKey generates a cryptographically secure key of a specified length with a given base as the seed.
// An error is returned if the length exceeds 64 or if there is an issue generating key components.
// Example result: `dP7P5LkjZfTjFVXCNcZg3iu1smwoHW` with `len` = 30
func GenKey(base, len int) ([]byte, error) {
	if len > 64 {
		return nil, fmt.Errorf("key length must be less than 64")
	}

	key := sha256.New()
	_, err := key.Write([]byte(strconv.Itoa(base) + strconv.Itoa(int(time.Now().Unix()))))
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	salt := make([]byte, 64)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	bytes := key.Sum(salt)
	encoded := base58.Encode(bytes)
	return []byte(encoded)[0:len], nil
}

func GenSepKey(base, parts int) ([]byte, error) {
	if parts > 64 {
		return nil, fmt.Errorf("key length must be less than 64")
	}

	key := sha256.New()
	_, err := key.Write([]byte(strconv.Itoa(base) + strconv.Itoa(int(time.Now().Unix()))))
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	salt := make([]byte, 64*lenOfPart)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	bytes := key.Sum(salt)
	encoded := base58.Encode(bytes)

	bytes = []byte(encoded)
	result := make([]byte, 0)
	partsCnt := 0

	for i := 0; i+1+lenOfPart < len(bytes) && partsCnt < parts; i += 1 + lenOfPart {
		result = append(result, bytes[i:i+lenOfPart]...)
		result = append(result, '-')
		partsCnt++
	}

	return result[:len(result)-1], nil
}
