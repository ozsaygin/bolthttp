package main

import (
	"fmt"
	"log"
)

type Server struct {
	address string
	port    int
}

func (s * Server) Serve() {
	log.Printf("Started to listen %s:%d", s.address, s.port)

}

func main() {
	// Example
	fmt.Println("Hello")
	server := &Server{"127.0.0.1", 6789}
	server.Serve()

}
