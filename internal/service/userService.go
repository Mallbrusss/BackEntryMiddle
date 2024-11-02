package service

import (

)

type UserService struct{
	jwt []byte
}

func NewUserService(jwt []byte) *UserService{
	return &UserService{
		jwt: jwt,
	}
}

func (us *UserService) Register(){

}

func (us *UserService) Authenticate(){

}

func (us *UserService) GetUserByUsername(){
	
}

func (us *UserService) ValidateToken(){

}