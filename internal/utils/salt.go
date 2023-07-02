package utils

import (
	"math/rand"
	"strings"
)

const (
	saltLength = 32
	letters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits     = "0123456789"
)

func GenerateRandomSalt(src rand.Source) (string, error) {
	var salt strings.Builder

	rnd := rand.New(src)

	for i := 0; i < saltLength; i++ {
		isDigit := rnd.Int()%2 == 0
		if isDigit {
			_, err := salt.WriteRune([]rune(digits)[rnd.Intn(len(digits))])
			if err != nil {
				return "", err
			}
		} else {
			smb := string([]rune(letters)[rnd.Intn(len(letters))])

			isLower := rnd.Int()%2 == 0
			if isLower {
				smb = strings.ToLower(smb)
			}

			_, err := salt.WriteString(smb)
			if err != nil {
				return "", err
			}
		}
	}

	return salt.String(), nil
}
