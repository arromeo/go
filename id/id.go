package id

import (
	"crypto/rand"
	"math/big"
)

func GenerateId() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const length = 8
	id := make([]byte, length)
	for i := range id {
		randomByte, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		id[i] = charset[randomByte.Int64()]
	}
	return string(id)
}
