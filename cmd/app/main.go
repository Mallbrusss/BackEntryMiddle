package main

import (
	"github.com/Mallbrusss/BackEntryMiddle/internal/server"
)

func init() {

}

func main() {

	serv := server.NewServer()
	serv.Run()
}
