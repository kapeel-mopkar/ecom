package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), plain); err != nil {
		log.Println("Login password did not match")
		return false
	}
	return true
}
