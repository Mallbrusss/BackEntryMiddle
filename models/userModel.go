package models

type User struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"pswd"`
}
