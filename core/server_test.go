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

func TestServer_Serve(t *testing.T) {
	type fields struct {
		address string
		port    int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"local", fields{"127.0.0.1", 8080}},
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
			buff := make([] byte, 4096)

			conn.Write([] byte("GET"))


			bufio.NewReader(conn).Read(buff)

			body := string(buff)

			fmt.Println("Body:" + body)
			if !strings.Contains(string(body), "ANSWER") {
				t.Errorf("Unexpected answer from server")
			}

			conn.Write([] byte("END"))
			log.Print("client sent end")
			buff = make([] byte, 4096)
			bufio.NewReader(conn).Read(
				buff)
			body = string(buff)
			if !strings.Contains(string(body), "ENDED") {
				t.Errorf("Ungracefull kill")
			}

		})
	}
}
