package internal

import (
	"math/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortName(min uint, max uint) string {
	length := rand.Intn(int(max)-int(min)+1) + int(min)

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(bytes)
}
