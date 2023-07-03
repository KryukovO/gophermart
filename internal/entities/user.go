package entities

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/KryukovO/gophermart/internal/utils"
)

var (
	ErrUserAlreadyExists    = errors.New("user with the same login already exists")
	ErrInvalidLoginPassword = errors.New("invalid login/password")
)

type User struct {
	ID                int64  `json:"-"`
	Login             string `json:"login"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"-"`
	Salt              string `json:"-"`
}

// Выполняет шифрование SHA-256 поля Password с добавлением соли.
func (user *User) Encrypt(secret []byte) error {
	if user.Salt == "" {
		salt, err := utils.GenerateRandomSalt(rand.NewSource(time.Now().UnixNano()))
		if err != nil {
			return err
		}

		user.Salt = salt
	}

	enc := hmac.New(sha256.New, secret)

	_, err := enc.Write([]byte(user.Password + user.Salt))
	if err != nil {
		return err
	}

	user.EncryptedPassword = hex.EncodeToString(enc.Sum(nil))

	return nil
}

// Возвращает ErrInvalidLoginPassword, если результат SHA-256 шифрования Password
// не соответствует EncryptedPassword с учетом Salt и secret.
// Всегда nil, если EncryptedPassword не установлен.
func (user *User) Validate(secret []byte) error {
	if user.EncryptedPassword == "" {
		return nil
	}

	enc := hmac.New(sha256.New, secret)

	_, err := enc.Write([]byte(user.Password + user.Salt))
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(enc.Sum(nil))

	if user.EncryptedPassword != hash {
		return ErrInvalidLoginPassword
	}

	return nil
}
