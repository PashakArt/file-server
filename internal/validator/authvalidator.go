package validator

import (
	"unicode"
	"unicode/utf8"

	"github.com/PashakArt/file-server/internal/domain"
)

const ()

func ValidateRegister(adminToken, token, login, password string) error {
	if token != adminToken {
		return domain.ErrTokenInvalid
	}

	return ValidateLogin(login, password)
}

func ValidateLogin(login, password string) error {
	if !IsValidLogin(login) {
		return domain.ErrLoginInvalid
	}

	if !IsValidPassword(password) {
		return domain.ErrPasswordInvalid
	}

	return nil
}

func IsValidLogin(login string) bool {
	if utf8.RuneCountInString(login) < 8 {
		return false
	}

	for _, ch := range login {
		if unicode.IsSpace(ch) {
			return false
		}

		if !isLatinLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func IsValidPassword(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}

	var (
		hasDigit  = false
		hasSymbol = false
		hasLower  = false
		hasUpper  = false
	)
	for _, ch := range password {
		if unicode.IsSpace(ch) {
			return false
		}

		switch {
		case isLatinLetter(ch) && unicode.IsLower(ch):
			hasLower = true
		case isLatinLetter(ch) && unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case !unicode.In(ch, unicode.Latin) && !unicode.IsDigit(ch):
			hasSymbol = true
		default:
			return false
		}
	}

	return hasDigit && hasSymbol && hasLower && hasUpper
}

func isLatinLetter(r rune) bool {
	return unicode.In(r, unicode.Latin) && unicode.IsLetter(r)
}
