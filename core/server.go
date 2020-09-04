/*
	Package httpproto implements an http server for HTTPv1.1.
*/
package core

import (
	"bufio"
	"encoding/json"
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

func mapPrettier(request map[string]string) string {
	b, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

func handleConnection(conn net.Conn) {

	// HTTP request format: [method] [resource] [http-version]\r\n
	//  \r\n states end of header.

	/* GET response

	HTTP/1.0 200 OK
	Etiam bibendum sapien ut est posuere pretium. Vestibulum a justo at sapien pharetra sagittis in eget lacus.

	HTTP/1.0 404 Not Found
	Sorry we don't have that file!

	HTTP/1.0 400 Bad Request
	Im sorry I just don't understand.

	*/

	//TODO: Implement the Content-Length header
	//TODO: Implement the Server header
	//TODO: Implement the Content-Type header
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	log.Printf("Connection established: %s", remoteAddr)

	// Wait for request from conn
	for {

		reader := bufio.NewReader(conn)
		buff := make([]byte, reader.Size())
		reader.Read(buff)
		lines := string(buff)
		request := make(map[string]string)
		dat := strings.Split(lines, "\n")

		for i, line := range dat {
			if i == 0 {

				header := strings.Split(line, " ")
				request["method"] = header[0]
				request["resource"] = header[1]
				request["version"] = header[2]
			}

			if strings.Contains(line, ":") {
				line := strings.ReplaceAll(line, "\u0000", "")
				pair := strings.Split(strings.TrimSpace(line), ":")
				request[pair[0]] = pair[1]
			}
		}

		log.Printf("Data received: %s", mapPrettier(request))
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
