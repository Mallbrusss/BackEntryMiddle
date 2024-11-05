package main

import (
	"log"

	"github.com/Mallbrusss/BackEntryMiddle/internal/server"
	"github.com/joho/godotenv"
)

func init() {

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading file: %v", err)
	}

	log.Println("Start Server")
	serv := server.NewServer()
	serv.Run()
}
