package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/mr-tron/base58"
	"strconv"
	"time"
)

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
