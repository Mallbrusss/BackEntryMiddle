package main

import(
	"internal/server"
)

func main(){
	serv := server.NewServer()
	serv.Run() 
}