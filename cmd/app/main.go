package main

import (
	"log"

	"github.com/0x0FACED/merch-shop/internal/server"
)

func main() {
	if err := server.StartHTTP(); err != nil {
		log.Fatalln("something went wrong: ", err)
	}
}
