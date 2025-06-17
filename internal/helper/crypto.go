package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	SESSION_ID_LENGTH = 16 // bytes
)

func Sha256(s string) string {
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return sha256Hash
}

func GenerateRandomBytesInHex(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Error while generating random bytes: %s", err)
	}
	randomBytesInHex := fmt.Sprintf("%x", buffer)
	return randomBytesInHex, nil
}

func GenerateSessionId() (string, error) {
	return GenerateRandomBytesInHex(SESSION_ID_LENGTH)
}

func IsSessionIdValid(sessionId string) bool {
	if len(sessionId) != SESSION_ID_LENGTH*2 {
		return false
	}
	_, err := hex.DecodeString(sessionId)
	return err == nil
}

func GenerateRandomAlphanumeric(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	if len(chars) != 62 {
		panic("Invalid total number of characters")
	}
	randomString := ""
	for len(randomString) < length {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			return "", fmt.Errorf("Error while generate a random byte: %s", err)
		}
		if b[0] < 248 { // 248 is 256 - (256 % 62) to avoid bias
			randomString += string(chars[b[0]%62])
		}
	}
	return randomString, nil
}

func HashPassword(password string, salt string) string {
	passwordWithSalt := fmt.Sprintf("%s:%s", salt, password)
	hashedPasswordWithSalt := Sha256(passwordWithSalt)
	hashedPasswordWithSaltAndHeader := fmt.Sprintf("sha256:%s:%s", salt, hashedPasswordWithSalt)
	return hashedPasswordWithSaltAndHeader
}

func IsPasswordCorrect(password string, hashedPasswordWithSaltAndHeader string) bool {
	parts := strings.Split(hashedPasswordWithSaltAndHeader, ":")
	if len(parts) != 3 {
		return false
	}
	if parts[0] != "sha256" {
		return false
	}
	salt := parts[1]
	if HashPassword(password, salt) != hashedPasswordWithSaltAndHeader {
		return false
	}
	return true
}
