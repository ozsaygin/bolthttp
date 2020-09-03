package core

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	address string
	port    int
}

func handleConnection(conn net.Conn) {

	// HTTP request format: [method] [resource] [http-version]\r\n
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	log.Printf("Connection established: %s", remoteAddr)

	// Wait for request from conn
	for {
		buffer := make([]byte, 4096)
		bufferSize, err := conn.Read(buffer)

		if bufferSize > 0 {
			log.Printf("Data received from connection %s: %s", conn.RemoteAddr(), string(buffer))

			data := strings.TrimSpace(string(buffer))
			log.Print(data)
			if strings.Contains(data, "GET") {
				fmt.Println("buffer is GET")
				conn.Write([]byte("ANSWER"))
				log.Print("Sent answer")
			}

			if  strings.Contains(data ,"END") {
				conn.Write([]byte("ENDED"))
				conn.Close()
				break
			}

			if err != nil {
				log.Printf("Cannot write to connection buffer")
			}
		}

		if err != nil {
			log.Printf("Something bad happened while reading data from %s\n", remoteAddr)
			log.Printf("Connection with %s closed...\n", remoteAddr)
			conn.Close()
			break
		}
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

		if err != nil {
			log.Printf("Connection from %s could not connect to server...", addr)
		}

		// Handle the connection in a separate go routine
		go handleConnection(conn)
	}

}
