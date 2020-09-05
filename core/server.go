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

func mapPrettier(request map[string]string) string {
	b, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// HTTP request format: [method] [resource] [http-version]\r\n
	//  \r\n states end of header.

	//TODO: Implement the Content-Length header
	//TODO: Implement the Server header
	//TODO: Implement the Content-Type header

	remoteAddr := conn.RemoteAddr()
	log.Printf("Connection established: %s", remoteAddr)
	fmt.Println("here132")
	// Wait for request from conn
	for {
		reader := bufio.NewReader(conn)
		buff := make([]byte, reader.Size())
		reader.Read(buff)

		data := string(buff)
		fmt.Println(data)
		request := make(map[string]string)

		lines := strings.Split(data, "\n")

		for i, line := range lines {

			// request headers have "\r\n" chars at the end of line
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

		fmt.Printf("Data received: %s", mapPrettier(request))

		// Process the request
		// Request dispatcher
		currentDir, err := os.Getwd()
		if err != nil {
			log.Println("Something bad happened while getting cwd")
		}


		resourceDir := "/www"
		switch request["method"] {
		case "GET":

			file, err := ioutil.ReadFile(currentDir + resourceDir  + request["resource"])
			if err != nil {


				conn.Write([]byte("HTTP/1.0 404 Not Found\r\nSorry we don't have that file!"))
				break

				log.Println("Cannot open the file")
				log.Println("[HTTP/1.0 404 Not Found]")
			} else {
				dataSent := file

				message := "HTTP/1.0 200 OK\""
				message += "\r\n"
				message += string(dataSent)
				log.Println(message)

				//color.Set(color.FgHiGreen)
				// TODO: Colorized server logs
				log.Println("[HTTP/1.0 200 OK]")
				break
			}


		case "POST":
			fmt.Println("POST method call")
		default:
			message := "HTTP/1.0 400 Bad Request"
			message += "\r\n"
			message += "Im sorry I just dont understand."
			msg := []byte(message)
			conn.Write(msg)

			log.Println("[HTTP/1.0 400 Bad Request]")

			if err != nil {
				log.Println("problematic close")
			}
			break
		}
	}
	log.Println("came out")
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
		log.Println("new connection")

		// Handle the connection in a separate go routine
		go handleConnection(conn)
	}

}
