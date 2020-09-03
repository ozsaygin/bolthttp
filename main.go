package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

type Server struct {
	address string
	port    int
}

func handleConnection(conn net.Conn) {

	// HTTP request rormat: [method] [resource] [http-version]\r\n

	log.Printf("Connection established: %s", conn.RemoteAddr())

	// Wait for request from conn
	for {
		data, _ := bufio.NewReader(conn).ReadByte()
		log.Printf("Data received from connection %s: %s", conn.RemoteAddr(), string(data))
	}
}

func (s *Server) Serve() {

	addr := s.address + ":" + strconv.Itoa(s.port)
	log.Printf("Started to listen %s...", addr)
	line, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Server could not listen %s...", addr)
	}

	// Start to accept incoming connections
	for {
		conn, err := line.Accept()
		defer conn.Close()

		if err != nil {
			log.Printf("Connection from %s could not connect to server...", conn.LocalAddr().String)
		}

		// Handle the connection in a seperate go routine
		go handleConnection(conn)
	}

}

func main() {

	// Example
	server := &Server{"127.0.0.1", 6789}
	server.Serve()

}
