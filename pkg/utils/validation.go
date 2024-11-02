package utils

import "regexp"

func IsValidLogin(login string) bool {
	//TODO: Мб перенести??
	loginRegex := regexp.MustCompile("^[a-zA-Z0-9]{8,}$`")
	return loginRegex.MatchString(login)
}

func IsValidPassword(password string) bool {
	//TODO: Мб перенести??
	passwordReg := regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*\\W).{8,}$")
	return passwordReg.MatchString(password)
}
