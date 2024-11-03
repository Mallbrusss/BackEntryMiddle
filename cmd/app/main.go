package main

import (
	"internal/server"
)

func init() {

}

func main() {

	serv := server.NewServer()
	serv.Run()
}
