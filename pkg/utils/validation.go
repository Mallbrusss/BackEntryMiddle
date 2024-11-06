package utils

import (
	"regexp"
)

// IsValidPassword проверяет логин на валидность
func IsValidLogin(login string) bool {
	loginRegex := regexp.MustCompile("^[a-zA-Zа-яА-Я0-9]{8,}$")
	return loginRegex.MatchString(login)
}

// IsValidPassword проверяет пароль на валидность
func IsValidPassword(password string) bool {

	switch {
	case len(password) < 8:
		return false
	case !regexp.MustCompile(`[a-zа-я]`).MatchString(password):
		return false
	case !regexp.MustCompile(`[A-ZА-Я]`).MatchString(password):
		return false
	case !regexp.MustCompile(`\d`).MatchString(password):
		return false
	case !regexp.MustCompile(`[^a-zA-Zа-яА-Я0-9]`).MatchString(password):
		return false
	default:
		return true
	}

}
