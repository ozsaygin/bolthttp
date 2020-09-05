package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &Server{
				Address: tt.fields.address,
				Port:    tt.fields.port,
			}
			// Start server
			go s.Serve()
			addr := s.Address + ":" + strconv.Itoa(s.Port)

			conn, err := net.Dial("tcp", addr)
			if err != nil {
				log.Println("Cannot connect to server")
			}
			msg := `GET /a.txt HTTP/1.1\n
User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)\n
Host: www.tutorialspoint.com\n
Accept-Language: en-us\n
Accept-Encoding: gzip, deflate\n
Connection: Keep-Alive\n`

			conn.Write([]byte(msg))

			buf,err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("handle me")
			}
			fmt.Println(buf)

			fmt.Println("Client Body:" + buf)
			conn.Close()
			fmt.Println("client closed")
		})
	}
}

// TODO: Write test for each method call [GET, POST, PUT, ...]
// TODO: Write test for found, missing resources and corrupted header