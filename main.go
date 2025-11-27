package main

import (
	"github.com/tomiok/dishwasher/internal/server"
	"log"
)

func main() {
	serv := server.New("7000")
	log.Fatal(serv.Start())
}
