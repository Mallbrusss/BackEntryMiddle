package utils

import "regexp"

func IsValidLogin(login string) bool {
	loginRegex := regexp.MustCompile("^[a-zA-Z0-9]`")
	return loginRegex.MatchString(login)
}

func IsValidPassword(password string) bool {
	passwordReg := regexp.MustCompile("^[a-zA-Z0-9]")
	return passwordReg.MatchString(password)
}
