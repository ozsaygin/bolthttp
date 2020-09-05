/*
	Package httpproto implements an http server for HTTPv1.1.
*/
package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	Address string
	Port    int
}

func mapPrettier(dictionary map[string]string) string {
	pretty, err := json.MarshalIndent(dictionary, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(pretty)
}

func handleConnection(conn net.Conn) {

	// TODO: Implement Content-Length header
	// TODO: Implement Server header
	// TODO: Implement Content-Type header

	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	log.Printf("Connection established: %s", remoteAddr)

	connected := true
	// Waits for requests from connection
	for connected {

		reader := bufio.NewReader(conn)
		buff := make([]byte, reader.Size())

		_, err := reader.Read(buff)

		if err != nil {
			log.Printf("Cannot read the buffer: %s", err)
		}

		data := string(buff)
		request := make(map[string]string)

		lines := strings.Split(data, "\n")

		if !(len(lines) >0) {
			continue
		}

		for i, line := range lines {

			// request headers have "\r\n" chars at the end of line
			if i == 0 {
				header := strings.Split(line, " ")
				request["method"] = header[0]
				request["resource"] = header[1]
				request["version"] =strings.ReplaceAll(header[2], "\r", "")

			}

			if strings.Contains(line, ":") {
				line := strings.ReplaceAll(line, "\u0000", "")
				pair := strings.Split(strings.TrimSpace(line), ":")
				request[pair[0]] = pair[1]
			}
		}

		fmt.Printf("Data received by server: \n%s", mapPrettier(request))

		// Process the request
		// Request dispatcher
		currentDir, err := os.Getwd()
		if err != nil {
			log.Println("Something bad happened while getting cwd")
		}

		// resourceDir := "/www"
		switch request["method"] {

		case "GET":

			// Read file requested
			file, err := ioutil.ReadFile(currentDir + request["resource"])

			// Resource not found
			if err != nil {

				conn.Write([]byte("HTTP/1.0 404 Not Found\r\n\nSorry we don't have that file!"))
				log.Println("HTTP/1.0 404 Not Found")
				conn.Close()
				connected = false

			} else {

				dataSent := file

				message := "HTTP/1.0 200 OK\r\n"
				message += "\n"
				message += string(dataSent)

				conn.Write([]byte(message))

				// TODO: Colorized server logs

				log.Println("HTTP/1.0 200 OK")

				// Panic occurs
				conn.Close()
				connected = false
			}

		case "POST":
			fmt.Println("POST method call")

		default:
			message := "HTTP/1.0 400 Bad Request\r\n"

			message += "\nIm sorry I just dont understand.\n"
			conn.Write([]byte(message))

			log.Println("HTTP/1.0 400 Bad Request")

			conn.Close()
			connected = false
		}
	}
}

func (s *Server) Serve() {

	addr := s.Address + ":" + strconv.Itoa(s.Port)
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
