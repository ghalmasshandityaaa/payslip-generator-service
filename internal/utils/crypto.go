package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(plaintext, encrypted string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(plaintext))
	return err == nil
}
