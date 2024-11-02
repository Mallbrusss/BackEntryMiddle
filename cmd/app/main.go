package main

import (
	"internal/server"
	"log"

	"github.com/joho/godotenv"
)

func init(){
	if err := godotenv.Load; err!= nil{
		log.Println(".env file not found")
	}
}

func main(){

	serv := server.NewServer()
	serv.Run() 
}