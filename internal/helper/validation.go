package helper

import (
	"net/mail"
	"regexp"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsPasswordValid(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLower && hasUpper && hasDigit
}

func IsStringNFKCNormalized(s string) bool {
	return norm.NFKC.IsNormalString(s)
}
