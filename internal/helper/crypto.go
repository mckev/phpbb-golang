package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
)

func Sha256(s string) string {
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return sha256Hash
}

func GenerateRandomSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return ""
	}
	saltString := fmt.Sprintf("%x", salt)
	return saltString
}

func HashPassword(password string, salt string) string {
	passwordWithSalt := fmt.Sprintf("%s:%s", salt, password)
	hashedPasswordWithSalt := Sha256(passwordWithSalt)
	return fmt.Sprintf("sha256:%s:%s", salt, hashedPasswordWithSalt)
}

func IsPasswordCorrect(password string, hashedPasswordWithSalt string) bool {
	parts := strings.Split(hashedPasswordWithSalt, ":")
	if len(parts) != 3 {
		return false
	}
	if parts[0] != "sha256" {
		return false
	}
	salt := parts[1]
	if HashPassword(password, salt) != hashedPasswordWithSalt {
		return false
	}
	return true
}
