package main

import (
	"internal/server"
	"log"

	"github.com/joho/godotenv"
)

func init() {

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env file not found")
	}
	
	serv := server.NewServer()
	serv.Run()
}
