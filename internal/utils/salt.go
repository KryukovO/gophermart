package utils

import (
	"math/rand"
	"strings"
)

const (
	saltLength = 32
	alphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomSalt(src rand.Source) (string, error) {
	var salt strings.Builder

	rnd := rand.New(src)
	runes := []rune(alphabet)

	for i := 0; i < saltLength; i++ {
		salt.WriteRune(runes[rnd.Intn(len(alphabet))])
	}

	return salt.String(), nil
}
