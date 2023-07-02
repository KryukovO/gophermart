package entities

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrUserAlreadyExists    = errors.New("user with the same login already exists")
	ErrInvalidLoginPassword = errors.New("invalid login/password")
)

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Salt     string `json:"-"`
}

// Выполняет шифрование SHA-256 поля Password с добавлением соли.
func (user *User) Encrypt(secret []byte) error {
	enc := hmac.New(sha256.New, secret)

	_, err := enc.Write([]byte(user.Password + user.Salt))
	if err != nil {
		return err
	}

	user.Password = hex.EncodeToString(enc.Sum(nil))

	return nil
}

// Проверка валидности password для пользователя
func (user *User) Validate(password string, secret []byte) error {
	enc := hmac.New(sha256.New, secret)

	_, err := enc.Write([]byte(password + user.Salt))
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(enc.Sum(nil))

	if user.Password != hash {
		return ErrInvalidLoginPassword
	}

	return nil
}
