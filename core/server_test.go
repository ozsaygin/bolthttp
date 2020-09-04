package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)
import "testing"

// TODO: Write a unit test for testing multi-client user handling

func TestServer_Serve(t *testing.T) {

	type fields struct {
		address string
		port    int
	}

	tests := []struct {
		name        string
		clientCount int
		fields      fields
	}{
		{"local", 1, fields{"127.0.0.1", 8080}},
		{"local", 4, fields{"127.0.0.1", 8081}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Server{
				address: tt.fields.address,
				port:    tt.fields.port,
			}
			// Start server
			go s.Serve()
			addr := s.address + ":" + strconv.Itoa(s.port)

			conn, err := net.Dial("tcp", addr)
			if err != nil {
				log.Println("Cannot connect to server")
			}
			buff := make([]byte, 4096)
			msg := `GET /hello.htm HTTP/1.1\n
User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)\n
Host: www.tutorialspoint.com\n
Accept-Language: en-us\n
Accept-Encoding: gzip, deflate\n
Connection: Keep-Alive\n`

			conn.Write([]byte(msg))

			bufio.NewReader(conn).Read(buff)

			body := string(buff)

			fmt.Println("Body:" + body)
			if !strings.Contains(string(body), "ANSWER") {
				t.Errorf("Unexpected answer from server")
			}

			conn.Write([]byte("END"))
			log.Print("client sent end")
			buff = make([]byte, 4096)
			bufio.NewReader(conn).Read(
				buff)
			body = string(buff)
			if !strings.Contains(string(body), "ENDED") {
				t.Errorf("Ungracefull kill")
			}

		})
	}
}
