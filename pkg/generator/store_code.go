package generator

import (
	"math/rand"
	"time"
)

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Generate(length int) string {
	time.Now().UnixNano()

	code := make([]byte, length)
	for i := range code {
		code[i] = base62[rand.Intn(len(base62))]
	}
	return string(code)
}
