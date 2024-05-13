package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("starting")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buffer := make([]byte, 1024)
		n, err := con.Read(buffer)
		if err != nil {
			fmt.Println("error reading response: ", err.Error())
			os.Exit(1)
		}

		requestContent := string(buffer[:n])
		requestLines := strings.Split(requestContent, "\r\n")
		route := strings.Split(requestLines[0], " ")[1]

		var response string
		switch route {
		case "/":
			response = "HTTP/1.1 200 OK\r\n\r\n"
		default:
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}

		_, err = con.Write([]byte(response))
		if err != nil {
			fmt.Println("error writing response: ", err.Error())
			os.Exit(1)
		}

		err = con.Close()
		if err != nil {
			fmt.Println("error closing connection: ", err.Error())
			os.Exit(1)
		}
	}
}
